---
name: add-swagger-api-docs
description: Add Swagger annotations for API handlers in manga-go using the correct route, response, and security conventions.

---

Add Swagger annotations for API handlers in manga-go.

**Resource:** ${input:resourceName:Resource name in PascalCase, for example Author, Genre, Comic}
**Route prefix:** ${input:routePrefix:Route path, for example /authors, /genres}
**Requires auth:** ${input:requireAuth:yes/no}

---

## Goal

Add complete Swagger comments for handler methods so `make swagger` can generate accurate docs aligned with real routes.

## Workflow

1. Locate the resource handler files in the route folder:
   - `internal/app/api/route/<resource>/create_<resource>.go`
   - `internal/app/api/route/<resource>/get_<resource>.go`
   - `internal/app/api/route/<resource>/get_<resource>s.go`
   - `internal/app/api/route/<resource>/get_all_<resource>s.go`
   - `internal/app/api/route/<resource>/update_<resource>.go`
   - `internal/app/api/route/<resource>/delete_<resource>.go`

2. Add a Swagger annotation block directly above each handler method.

3. Ensure annotations exactly match endpoints in `route.go`:
   - `@Router` must use the exact path and lowercase method (`[get]`, `[post]`, `[put]`, `[delete]`, `[patch]`)
   - `@Tags` must use the resource name in PascalCase
   - `@Param` must use the correct type: `path`, `query`, `body`, `formData`

4. Response conventions:
   - GET list (paginated): `@Success 200 {object} response.Result`
   - GET detail: `@Success 200 {object} response.Result`
   - POST/PUT/DELETE: `@Success 200 {object} response.Result`
   - Common failures:
     - `@Failure 400 {object} response.Result`
     - `@Failure 401 {object} response.Result` (if endpoint requires auth)
     - `@Failure 500 {object} response.Result`

5. Security rule:
   - Add `@Security AccessToken` for authenticated endpoints.
   - Do not add it for signup/signin/reset-password.

6. If the endpoint uploads a file:
   - Add `@Accept multipart/form-data`
   - Add `@Param file formData file true "File to upload"`

7. Regenerate Swagger docs:
   ```bash
   make swagger
   ```

8. Quick verification:
   - Confirm docs were updated: `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`
   - Run build:
   ```bash
   go build ./...
   ```

## Annotation example

```go
// @Summary      Create author
// @Description  Create a new author in the system
// @Tags         Author
// @Accept       json
// @Produce      json
// @Param        body  body      authorrequest.CreateAuthorRequest  true  "Author creation request"
// @Success      200   {object}  response.Result
// @Failure      400   {object}  response.Result
// @Failure      401   {object}  response.Result
// @Failure      500   {object}  response.Result
// @Security     AccessToken
// @Router       /authors [post]
func (h *AuthorHandler) createAuthor(c *gin.Context) {
    // ...
}
```

## Checklist

- [ ] All resource handler methods have Swagger annotations
- [ ] `@Router` matches route.go exactly
- [ ] Security annotations match auth requirements
- [ ] `make swagger` runs successfully
- [ ] `go build ./...` runs successfully
