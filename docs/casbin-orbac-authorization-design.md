# Casbin OrBAC Authorization Design

## 1. Current State

The project already has the first pieces for authorization:

- `internal/pkg/casbin` creates a Casbin enforcer with GORM adapter and Redis watcher.
- `internal/pkg/casbin/model.conf` defines the OrBAC-style request/policy/grouping model.
- `roles`, `permissions`, `users_roles`, and `roles_permissions` tables already exist.
- Seeders create coarse permissions such as `comic:read`, `chapter:write`, `role:manage`.
- Routes use `AuthMiddleware.RequireJwt` plus `AuthzMiddleware.Require(action, resource, resolvers...)` where authorization is required.
- Domain ownership data already exists:
  - `users.translation_group_id`
  - `translation_groups.owner_id`
  - `comics.uploaded_by_id`
  - `comics.translation_group_id`
  - `chapters.uploaded_by_id`
  - comments, ratings, reading histories, notifications have user ownership fields.

Casbin policies are persisted in `casbin_rule` by the GORM Adapter. Runtime policy changes are propagated early through the Redis Watcher.

## 2. Goal

Use Casbin as the decision engine and express authorization in an OrBAC-style model:

- Organization: the authorization scope, usually `platform` or a translation group.
- Role: `admin`, `reader`, `translator`, `group_owner`, `group_member`, `moderator`.
- Activity: business action such as `read`, `create`, `update`, `delete`, `publish`, `manage`.
- View: resource class such as `comic`, `chapter`, `translation_group`, `comment`, `rating`.
- Context: condition such as `any`, `owner`, `group_member`, `group_owner`, `published`.

The practical result should be:

- Middleware checks for simple permissions and object-sensitive permissions such as owner/group/published.
- Routes explicitly pass `action` and `resource` to middleware.
- Resource resolvers are optional and lazy: middleware checks `any` first and only queries the database when object context is needed.
- Existing role and permission APIs remain admin-facing and update Casbin through `PolicyManager`, using Casbin API calls backed by GORM Adapter.

## 3. Recommended Model

Use Casbin RBAC with domains to represent OrBAC organization scope.

Request:

```ini
r = sub, org, act, obj, ctx
```

Policy:

```ini
p = org, role, act, view, ctx, eft
```

Role assignment:

```ini
g = sub, role, org
```

Matcher:

```ini
m = (g(r.sub, p.role, r.org) || g(r.sub, p.role, p.org)) &&
    (p.org == r.org || p.org == "platform") &&
    (p.view == "*" || keyMatch(r.obj, p.view)) &&
    (p.act == r.act || p.act == "manage") &&
    (p.ctx == r.ctx || p.ctx == "any") &&
    (p.eft == "allow")
```

Suggested `internal/pkg/casbin/model.conf`:

```ini
[request_definition]
r = sub, org, act, obj, ctx

[policy_definition]
p = org, role, act, view, ctx, eft

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (g(r.sub, p.role, r.org) || g(r.sub, p.role, p.org)) && (p.org == r.org || p.org == "platform") && (p.view == "*" || keyMatch(r.obj, p.view)) && (p.act == r.act || p.act == "manage") && (p.ctx == r.ctx || p.ctx == "any")
```

Notes:

- `org` is a domain. Use `platform` for global authorization.
- For group-specific authorization, use `tg:<uuid>` as org.
- The matcher checks both `r.org` and `p.org` role assignments so a platform role such as `admin` can authorize requests against a translation-group org.
- `obj` should normally be a normalized resource class such as `comic`, `chapter`, or `translation_group`. If instance-level patterns are needed later, use strings such as `comic/<uuid>` and policies such as `comic/*`.
- `view` can be exact or pattern-based, for example `comic`, `chapter`, `comment/*`, or `*`.
- `ctx` is computed by application code, not by Casbin. Casbin only compares the computed context with policy.

## 4. Policy Examples

Global platform roles:

```csv
g, <admin-user-id>, admin, platform
g, <reader-user-id>, reader, platform
g, <translator-user-id>, translator, platform

p, platform, admin, manage, *, any, allow
p, platform, reader, read, comic, published, allow
p, platform, reader, read, chapter, published, allow
p, platform, reader, create, comment, any, allow
p, platform, reader, update, comment, owner, allow
p, platform, reader, delete, comment, owner, allow
p, platform, reader, create, rating, any, allow
p, platform, reader, update, rating, owner, allow
p, platform, reader, delete, rating, owner, allow

p, platform, translator, create, comic, any, allow
p, platform, translator, update, comic, owner, allow
p, platform, translator, update, comic, group_member, allow
p, platform, translator, create, chapter, group_member, allow
p, platform, translator, update, chapter, owner, allow
p, platform, translator, update, chapter, group_member, allow
p, platform, translator, publish, chapter, group_member, allow
```

Translation group roles:

```csv
g, <owner-user-id>, group_owner, tg:<group-id>
g, <member-user-id>, group_member, tg:<group-id>

p, tg:<group-id>, group_owner, manage, translation_group, group_owner, allow
p, tg:<group-id>, group_owner, manage, comic, group_member, allow
p, tg:<group-id>, group_owner, manage, chapter, group_member, allow
p, tg:<group-id>, group_member, create, chapter, group_member, allow
p, tg:<group-id>, group_member, update, chapter, owner, allow
```

## 5. Context Computation

Casbin should not query the database. The app should load resource facts and compute one or more contexts before calling `Enforce`.

Suggested context values:

- `any`: no resource relationship needed.
- `owner`: current user owns the object.
- `group_member`: current user is in the object's translation group.
- `group_owner`: current user owns the object's translation group.
- `published`: object is public/readable.
- `self`: object is the current user.

For resources in this codebase:

| Resource | Owner/context source |
| --- | --- |
| user | `users.id == currentUser.ID` gives `self` |
| translation_group | `translation_groups.owner_id == currentUser.ID` gives `group_owner` |
| comic | `uploaded_by_id == currentUser.ID` gives `owner`; `translation_group_id == currentUser.translation_group_id` gives `group_member`; `is_published` gives `published` |
| chapter | `uploaded_by_id == currentUser.ID` gives `owner`; parent comic's group gives `group_member`; `is_published` gives `published` |
| comment | `comments.user_id == currentUser.ID` gives `owner` |
| rating | `ratings.user_id == currentUser.ID` gives `owner` |
| reading_history | `reading_histories.user_id == currentUser.ID` gives `owner` |
| notification | `user_notifications.user_id == currentUser.ID` gives `owner` |

If multiple contexts are true, check them in priority order:

```go
contexts := []string{"owner", "group_owner", "group_member", "published", "any"}
for _, ctx := range contexts {
    ok, err := enforcer.Enforce(userID, org, action, object, ctx)
    if err != nil {
        return err
    }
    if ok {
        return nil
    }
}
return ErrForbidden
```

## 6. Authorization Service

Add a small application-level service around Casbin. Do not call Casbin directly from every handler.

Suggested package:

```text
internal/pkg/authorization
```

Suggested API:

```go
type Authorizer struct {
    enforcer *casbin.Enforcer
}

type Request struct {
    Subject string
    Org     string
    Action  string
    Object  string
    Context string
}

func (a *Authorizer) Enforce(ctx context.Context, req Request) error
func (a *Authorizer) EnforceAny(ctx context.Context, req Request, contexts []string) error
```

Also add helpers:

```go
func PlatformOrg() string { return "platform" }
func TranslationGroupOrg(id uuid.UUID) string { return "tg:" + id.String() }
func Subject(id uuid.UUID) string { return id.String() }
```

## 7. Middleware Design

Authorization checks are middleware-owned. Routes pass the business action, resource class, and optional lazy resolvers:

```go
authzmiddleware.Require(m, authorization.ActionUpdate, authorization.ObjectComic, m.Comic())
```

The middleware flow is:

1. Enforce `action/resource/any` against the `platform` org.
2. If allowed, continue without loading the target resource.
3. If denied and no resolver is configured, return `403`.
4. If denied and resolvers exist, run resolvers one by one and enforce computed contexts such as `owner`, `group_member`, `group_owner`, `published`, or `self`.
5. Stop as soon as one resolver context is allowed, so later resolvers do not query the database unnecessarily.

Static examples:

- `POST /roles` -> `manage role`
- `POST /permissions` -> `manage permission`
- `POST /tags` -> `create tag`
- `DELETE /genres/:slug` -> `delete genre`
- `POST /files/upload` -> `create file`

Object-sensitive examples:

- `PUT /comics/:slug` -> `update comic` with `Comic()` resolver
- `PATCH /comics/:slug/chapters/:chapterSlug/publish` -> `publish chapter` with `Chapter()` then `ComicGroupFromContext()`
- `PUT /comments/:id` -> `update comment` with `CommentParam("id")`
- `PATCH /users/:id` -> `update user` with `UserParam("id")`

## 8. Policy Persistence

Do not add a custom runtime policy syncer. Policy writes go through Casbin APIs:

```go
enforcer.AddPolicy(...)
enforcer.RemovePolicy(...)
enforcer.AddGroupingPolicy(...)
enforcer.RemoveGroupingPolicy(...)
```

The GORM Adapter automatically persists these changes to `casbin_rule`. Redis Watcher reloads policy quickly in other application instances.

`PolicyManager` is only a thin application wrapper around these Casbin calls. It is used by role, permission, user-role, and translation-group ownership flows.

## 9. Route Permission Matrix

Suggested initial mapping:

| Endpoint pattern | Action | Object | Context |
| --- | --- | --- | --- |
| `GET /comics`, `GET /comics/:slug` | read | comic | published/group_member/owner |
| `POST /comics` | create | comic | any |
| `PUT /comics/:slug` | update | comic | owner/group_member |
| `PATCH /comics/:slug/status` | update | comic | owner/group_member |
| `PATCH /comics/:slug/publish` | publish | comic | group_owner/admin |
| `DELETE /comics/:slug` | delete | comic | owner/group_owner/admin |
| `GET /comics/:slug/chapters` | read | chapter | published/group_member/owner |
| `POST /comics/:slug/chapters` | create | chapter | group_member |
| `PUT /comics/:slug/chapters/:chapterSlug` | update | chapter | owner/group_member |
| `PUT /comics/:slug/chapters/:chapterSlug/pages` | update | chapter | owner/group_member |
| `PATCH /comics/:slug/chapters/:chapterSlug/publish` | publish | chapter | group_member/group_owner |
| `POST /comments` | create | comment | any |
| `PUT /comments/:id` | update | comment | owner |
| `DELETE /comments/:id` | delete | comment | owner/moderator/admin |
| `POST /ratings` | create | rating | any |
| `PUT /ratings/:id` | update | rating | owner |
| `DELETE /ratings/:id` | delete | rating | owner/moderator/admin |
| `/roles/**` | manage | role | any |
| `/permissions/**` | manage | permission | any |
| `/translation-groups/:slug/**` | manage | translation_group | group_owner |

## 10. Data Migration Strategy

The migration creates `casbin_rule`, inserts built-in policies, and backfills once from existing tables:

- `users_roles` -> `g, <user_id>, <role.name>, platform`
- `roles_permissions` -> `p, platform, <role.name>, <action>, <resource>, any, allow`
- `users.translation_group_id` -> `g, <user_id>, group_member, tg:<translation_group_id>`
- `translation_groups.owner_id` -> `g, <owner_id>, group_owner, tg:<translation_group_id>`

## 11. Policy Sync Rules

Current permission names are `resource:action`, for example `comic:read`. Convert them into Casbin policies:

```text
permission.name = "<resource>:<action>"
Casbin p = platform, <role.name>, <action>, <resource>, any, allow
```

For user roles:

```text
users_roles row -> g, <user_id>, <role.name>, platform
```

For translation group membership:

```text
users.translation_group_id -> g, <user_id>, group_member, tg:<translation_group_id>
translation_groups.owner_id -> g, <owner_id>, group_owner, tg:<translation_group_id>
```

When these change, update Casbin through `PolicyManager`:

- `UserService.AssignRoles`
- `UserService.RemoveRole`
- `RoleService.AssignPermissions`
- `RoleService.RemovePermission`
- `TranslationGroupService.CreateTranslationGroup`
- `TranslationGroupService.TransferOwnership`
- any future group member join/leave endpoint

There is no startup reconciliation job. The one-time migration backfills existing data, and all future writes update Casbin directly.

## 12. Implementation Plan

1. Fill `internal/pkg/casbin/model.conf` with the OrBAC-style model.
2. Register `casbin.NewEnforcer` in `internal/app/fx.go`.
3. Add `internal/pkg/authorization` with `Authorizer` and helpers.
4. Add `internal/app/middleware/authz`.
5. Add `PolicyManager` as a thin wrapper over Casbin management APIs.
6. Apply middleware checks for admin/catalog endpoints.
7. Apply middleware checks with lazy resolvers for comics, chapters, groups, comments, ratings, histories, and user self-update.
8. Add tests:
   - Casbin model tests for allow/deny examples.
   - Authorizer tests for context priority.
   - Service tests for owner/group forbidden cases.
   - Route tests for admin-only endpoints.

## 13. Important Security Fixes To Include

While adding authorization, also tighten these existing flows:

- `CommentService.UpdateComment` and `DeleteComment` currently do not verify the current user owns the comment.
- `ComicService.UpdateComic`, `PublishComic`, `UpdateComicStatus`, and `DeleteComic` do not verify uploader or group membership.
- `ChapterService.UpdateChapter`, `UpdateChapterPages`, and `PublishChapter` do not verify uploader or group membership.
- `TranslationGroupService.UpdateTranslationGroup`, `DeleteTranslationGroup`, `TransferOwnership`, and logo update should require group owner or platform admin.
- User role management routes should require `role manage`.
- Permission management routes should require `permission manage`, not merely JWT.

## 14. Suggested First Rollout

Start with a conservative policy:

- `admin` can manage everything.
- `reader` can read published comics/chapters and manage own comments, ratings, histories, and notifications.
- `translator` can create comics and chapters, and update/publish resources belonging to their translation group.
- `group_owner` can manage its own translation group, comics, and chapters.

Do not expose deny policies in the first rollout unless there is a clear need. Keep `eft = allow` only, because the current product rules are simpler as positive permissions plus object context.
