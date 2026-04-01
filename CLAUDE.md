# manga-go — Hướng dẫn kiến trúc & Convention

## Tổng quan dự án

REST API phục vụ ứng dụng manga, xây dựng bằng Go với stack:

| Thành phần | Thư viện |
|---|---|
| HTTP framework | `gin-gonic/gin` |
| Dependency Injection | `uber/fx` |
| ORM | `gorm.io/gorm` + driver PostgreSQL |
| Cache / Queue backend | Redis (`go-redis/v9`) |
| Task queue | `hibiken/asynq` |
| Migrations | `rubenv/sql-migrate` (CLI `sql-migrate`) |
| Logger | `uber/zap` (wrapper ở `internal/pkg/logger`) |
| Config | `spf13/viper` |
| Tracing | OpenTelemetry (`go.opentelemetry.io/otel`) |
| Auth | JWT cookie (`golang-jwt/jwt/v5`) + Redis blacklist |

---

## Cấu trúc thư mục

```
manga-go/
├── cmd/
│   ├── dev/main.go          # Entrypoint API server
│   └── queue/main.go        # Entrypoint async worker
├── configs/
│   └── dbconfig.yml         # sql-migrate config
├── internal/
│   ├── app/
│   │   ├── fx.go            # Root fx.Module
│   │   ├── api/
│   │   │   ├── fx.go
│   │   │   ├── common/
│   │   │   │   ├── route.go              # Route interface + ProvideAsRoute
│   │   │   │   └── response/response.go  # Response types & helpers
│   │   │   ├── route/
│   │   │   │   ├── fx.go                 # Tổng hợp route modules
│   │   │   │   ├── author/               # Author CRUD routes
│   │   │   │   ├── genre/                # Genre CRUD routes
│   │   │   │   └── user/                 # User auth routes
│   │   │   └── server/                   # Gin setup, HTTP server, lifecycle
│   │   └── middleware/auth/              # JWT auth middleware
│   └── pkg/
│       ├── common/
│       │   ├── gorm.go      # SqlModel (base model)
│       │   └── paging.go    # Paging struct
│       ├── model/           # GORM models
│       ├── repo/
│       │   ├── base/repo.go # Generic BaseRepository[T]
│       │   └── <resource>/  # Concrete repos
│       ├── request/<resource>/ # Request DTOs
│       ├── services/<resource>/ # Business logic
│       └── ...
├── migrations/              # SQL migration files
└── CLAUDE.md
```

---

## Quy ước & Pattern bắt buộc

### 1. Dependency Injection (uber/fx)

Mỗi package phải có file `fx.go` khai báo `var Module = fx.Module(...)`.

**Params pattern** — mọi constructor nhận một struct với `fx.In`:
```go
type FooServiceParams struct {
    fx.In
    Logger  *logger.Logger
    FooRepo *foorepo.FooRepo
}

func NewFooService(p FooServiceParams) *FooService {
    return &FooService{logger: p.Logger, fooRepo: p.FooRepo}
}
```

**Module registration** — module con phải được thêm vào module cha:
- Route mới → `internal/app/api/route/fx.go`
- Repo mới  → `internal/pkg/repo/fx.go`
- Service mới → `internal/pkg/services/fx.go`

### 2. Model

Mọi model embed `common.SqlModel` (UUID PK + timestamps + soft-delete):
```go
type Foo struct {
    common.SqlModel
    Name string `json:"name" gorm:"column:name"`
}
func (Foo) TableName() string { return "foos" }
```

File: `internal/pkg/model/<resource>.go`

### 3. Repository

Embed `BaseRepository[T]` — không cần implement lại CRUD:
```go
type FooRepo struct {
    *base.BaseRepository[model.Foo]
}
func NewFooRepo(db *gorm.DB) *FooRepo {
    return &FooRepo{BaseRepository: &base.BaseRepository[model.Foo]{DB: db}}
}
```

`BaseRepository[T]` cung cấp sẵn:
- `FindOne`, `FindAll`, `FindPaginated`
- `Create`, `CreateList`
- `Update`, `UpdateWithTransaction`
- `DeleteSoft`, `DeletePermanently`
- `Upsert`, `UpsertMany`
- `CountAll`

Điều kiện WHERE truyền dưới dạng `[]any` với `clause.Eq`, `clause.Where`, hoặc `common.JoinExpr`.

File: `internal/pkg/repo/<resource>/repo.go` + `fx.go`

### 4. Request DTO

```go
type CreateFooRequest struct {
    Name string `json:"name" binding:"required"`
}
type UpdateFooRequest struct {
    Name string `json:"name" binding:"required"`
}
```

File: `internal/pkg/request/<resource>/create.go`, `update.go`

### 5. Service

Mỗi operation nằm trong file riêng, trả về `response.Result`:
```go
// internal/pkg/services/foo/create.go
func (s *FooService) CreateFoo(ctx context.Context, req *foorequest.CreateFooRequest) response.Result {
    foo := model.Foo{Name: req.Name}
    if err := s.fooRepo.Create(ctx, &foo); err != nil {
        s.logger.Error("Failed to create foo", "error", err)
        return response.ResultErrDb(err)
    }
    return response.ResultSuccess("Foo created successfully", foo)
}
```

Kiểm tra record not found:
```go
import "errors"
import "gorm.io/gorm"
if errors.Is(err, gorm.ErrRecordNotFound) {
    return response.ResultNotFound("Foo")
}
return response.ResultErrDb(err)
```

Files: `service.go`, `fx.go`, `create.go`, `get.go`, `list.go`, `list_all.go`, `update.go`, `delete.go`

### 6. API Route

4 files bắt buộc cho mỗi resource:

**`handler.go`** — struct handler + params:
```go
type FooHandler struct{ fooService *fooservice.FooService }
type FooHandlerParams struct {
    fx.In
    FooService *fooservice.FooService
}
func NewFooHandler(p FooHandlerParams) *FooHandler {
    return &FooHandler{fooService: p.FooService}
}
```

**`route.go`** — đăng ký endpoints:
```go
func (r *FooRoute) Setup() {
    rg := r.r.Group("/foos", r.authMiddleware.RequireJwt)
    rg.GET("/", r.fooHandler.getFoos)
    rg.GET("/:id", r.fooHandler.getFoo)
    rg.POST("/", r.fooHandler.createFoo)
    rg.PUT("/:id", r.fooHandler.updateFoo)
    rg.DELETE("/:id", r.fooHandler.deleteFoo)
}
```

**`fx.go`**:
```go
var Module = fx.Module(
    "foo-route",
    common.ProvideAsRoute(NewFooRoute),
    fx.Provide(NewFooHandler),
)
```

**Handler methods** (`create_foo.go`, `get_foo.go`, ...):
```go
func (h *FooHandler) createFoo(c *gin.Context) {
    var req foorequest.CreateFooRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ResultInvalidRequestErr(err).ResponseResult(c)
        return
    }
    result := h.fooService.CreateFoo(c.Request.Context(), &req)
    result.ResponseResult(c)
}
```

Lấy UUID từ path param:
```go
id, err := uuid.Parse(c.Param("id"))
if err != nil {
    response.ResultError("Invalid id").ResponseResult(c)
    return
}
```

Lấy paging từ query:
```go
var paging common.Paging
if err := c.ShouldBindQuery(&paging); err != nil {
    response.ResultInvalidRequestErr(err).ResponseResult(c)
    return
}
```

### 7. Migration

Tên file: `YYYYMMDD_HHMMSS_<mô_tả>.sql`

```sql
-- +migrate Up
CREATE TABLE foos (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER update_foos_updated_at
BEFORE UPDATE ON foos
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_foos_updated_at ON foos;
DROP TABLE foos;
```

### 8. Response helpers

| Function | Dùng khi |
|---|---|
| `response.ResultSuccess(msg, data)` | Thành công |
| `response.ResultNotFound(entity)` | Record không tồn tại |
| `response.ResultErrDb(err)` | Lỗi DB (500) |
| `response.ResultErrInternal(err)` | Lỗi server khác (500) |
| `response.ResultError(msg)` | Lỗi domain (400) |
| `response.ResultInvalidRequestErr(err)` | Lỗi bind/validate request (400) |
| `response.ResultUnauthorized()` | Chưa xác thực (401) |
| `response.ResponsePaginationData(elements, total)` | Wrap dữ liệu phân trang |

### 9. Swagger/OpenAPI Documentation

Mỗi handler method phải có comments Swagger phía trên function signature:

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
func (h *AuthorHandler) createAuthor(c *gin.Context) { ... }
```

**Quy ước annotation:**
- `@Tags`: Sử dụng tên resource dạng PascalCase (User, Author, Genre, Comic, Chapter, TranslationGroup, Role, Permission, File)
- `@Param`: Loại `path` cho URL param, `query` cho query string, `body` cho JSON body, `formData` cho multipart
- `@Success`: Trả về `response.Result` cho single entity hoặc API responses
- `@Security`: Thêm `AccessToken` cho endpoints cần xác thực. KHÔNG dùng cho endpoint signup, signin, reset-password
- `@Router`: Đường dẫn chính xác từ route.go, method lowercase ([get], [post], [put], [delete], [patch])

**Loại endpoint:**
- GET list (paginated): `@Success 200 {object} response.PaginationResponse`
- GET detail: `@Success 200 {object} response.Result`
- POST/PUT/DELETE: `@Success 200 {object} response.Result`
- File upload: `@Accept multipart/form-data` + `@Param file formData file true "File to upload"`

Sau khi thêm annotations, chạy:
```bash
make swagger
# hoặc
swag init -g cmd/dev/main.go -o docs/ --parseDependency --parseInternal
```

Các file docs/docs.go, swagger.json, swagger.yaml được generate automatically.

---

## Checklist khi thêm resource mới

```
[ ] migrations/YYYYMMDD_HHMMSS_create_<resource>_table.sql
[ ] internal/pkg/model/<resource>.go
[ ] internal/pkg/request/<resource>/create.go
[ ] internal/pkg/request/<resource>/update.go
[ ] internal/pkg/repo/<resource>/repo.go
[ ] internal/pkg/repo/<resource>/fx.go
[ ] internal/pkg/repo/fx.go                    ← thêm module
[ ] internal/pkg/services/<resource>/service.go
[ ] internal/pkg/services/<resource>/fx.go
[ ] internal/pkg/services/<resource>/create.go
[ ] internal/pkg/services/<resource>/get.go
[ ] internal/pkg/services/<resource>/list.go
[ ] internal/pkg/services/<resource>/list_all.go
[ ] internal/pkg/services/<resource>/update.go
[ ] internal/pkg/services/<resource>/delete.go
[ ] internal/pkg/services/fx.go                ← thêm module
[ ] internal/app/api/route/<resource>/handler.go
[ ] internal/app/api/route/<resource>/route.go
[ ] internal/app/api/route/<resource>/fx.go
[ ] internal/app/api/route/<resource>/create_<resource>.go
[ ] internal/app/api/route/<resource>/get_<resource>.go
[ ] internal/app/api/route/<resource>/get_<resource>s.go
[ ] internal/app/api/route/<resource>/get_all_<resource>s.go
[ ] internal/app/api/route/<resource>/update_<resource>.go
[ ] internal/app/api/route/<resource>/delete_<resource>.go
[ ] internal/app/api/route/fx.go               ← thêm module
```

---

## Lệnh thường dùng

```bash
# Chạy dev server (live reload)
air -c api.air.toml

# Build nhanh
go build ./...

# Generate Swagger documentation
make swagger

# Chạy migration lên
sql-migrate up -config=configs/dbconfig.yml

# Rollback migration
sql-migrate down -config=configs/dbconfig.yml

# Kiểm tra migration status
sql-migrate status -config=configs/dbconfig.yml

# Chạy queue worker
go run cmd/queue/main.go
```

---

## Naming conventions

| Loại | Pattern | Ví dụ |
|---|---|---|
| Package repo | `<resource>repo` | `authorRepo`, `genreRepo` |
| Package service | `<resource>service` | `authorservice`, `genreservice` |
| Package route | `<resource>route` | `authorroute`, `genreroute` |
| Package request | `<resource>request` | `authorrequest`, `genrerequest` |
| fx.Module name | `"<resource>-route"` | `"author-route"`, `"genre-route"` |
| DB table | snake_case plural | `authors`, `genres`, `users` |
| Go struct | PascalCase singular | `Author`, `Genre`, `User` |
| Route group | `/snake_case_plural` | `/authors`, `/genres` |
