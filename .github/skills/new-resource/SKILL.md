---
name: new-resource
description: 'Tạo đầy đủ một CRUD resource mới cho dự án manga-go theo đúng kiến trúc, convention fx module, repository, service, route, migration, và build verification. Use when: scaffold resource mới, thêm CRUD module mới, tạo migration/model/repo/service/route đồng bộ.'
argument-hint: 'Tên resource dạng singular PascalCase, ví dụ: Tag, Category, Publisher'
---

# New Resource

Tạo đầy đủ một CRUD resource mới cho dự án manga-go.

## Khi nào dùng

- Khi cần scaffold một resource CRUD mới theo đúng kiến trúc dự án.
- Khi cần tạo đồng bộ migration, model, request DTO, repository, service, route và đăng ký fx module.
- Khi cần đảm bảo resource mới build được với `go build ./...`.

## Đầu vào bắt buộc

- Tên resource ở dạng singular PascalCase, ví dụ: `Tag`, `Category`, `Publisher`.

Nếu người dùng chưa cung cấp tên resource hoặc chưa mô tả các field domain cần có, hỏi ngắn gọn để lấy đủ thông tin trước khi tạo file.

## Quy trình thực hiện

Đầu tiên, đọc `CLAUDE.md` và `./../copilot-instructions.md` nếu cần để nắm convention của dự án. Sau đó tạo toàn bộ file theo đúng thứ tự dưới đây.

### Bước 1 - Migration

Tạo file:

`migrations/<YYYYMMDD>_<HHMMSS>_create_<resource>_table.sql`

Lấy timestamp thực tế bằng lệnh:

```bash
date +%Y%m%d_%H%M%S
```

Template migration:

```sql
-- +migrate Up
CREATE TABLE <resource>s (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    -- them cac cot domain tai day
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TRIGGER update_<resource>s_updated_at
BEFORE UPDATE ON <resource>s
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_<resource>s_updated_at ON <resource>s;
DROP TABLE <resource>s;
```

### Bước 2 - Model

Tạo file `internal/pkg/model/<resource>.go`.

Yêu cầu:

- Embed `common.SqlModel`.
- Khai báo struct model dạng singular PascalCase.
- Implement `TableName() string` trả về tên bảng dạng snake_case plural.

### Bước 3 - Request DTOs

Tạo các file:

- `internal/pkg/request/<resource>/create.go`
- `internal/pkg/request/<resource>/update.go`

Yêu cầu:

- Dùng tên struct `Create<Resource>Request` và `Update<Resource>Request`.
- Dùng tag `binding:"required"` cho các trường bắt buộc.

### Bước 4 - Repository

Tạo các file:

- `internal/pkg/repo/<resource>/repo.go`
- `internal/pkg/repo/<resource>/fx.go`

Yêu cầu:

- Repository embed `*base.BaseRepository[model.<Resource>]`.
- Tạo constructor `New<Resource>Repo(db *gorm.DB) *<Resource>Repo`.
- Khai báo `var Module = fx.Module(...)` trong `fx.go`.
- Cập nhật `internal/pkg/repo/fx.go` để thêm module mới.

### Bước 5 - Service

Tạo các file:

- `internal/pkg/services/<resource>/service.go`
- `internal/pkg/services/<resource>/fx.go`
- `internal/pkg/services/<resource>/create.go`
- `internal/pkg/services/<resource>/get.go`
- `internal/pkg/services/<resource>/list.go`
- `internal/pkg/services/<resource>/update.go`
- `internal/pkg/services/<resource>/delete.go`

Yêu cầu:

- Constructor dùng params struct với `fx.In`.
- `get.go` phải dùng `errors.Is(err, gorm.ErrRecordNotFound)` để trả `ResultNotFound`.
- `list.go` dùng `FindPaginated` và `response.ResponsePaginationData`.
- Cập nhật `internal/pkg/services/fx.go` để thêm module mới.

### Bước 6 - API Route

Tạo các file:

- `internal/app/api/route/<resource>/handler.go`
- `internal/app/api/route/<resource>/route.go`
- `internal/app/api/route/<resource>/fx.go`
- `internal/app/api/route/<resource>/create_<resource>.go`
- `internal/app/api/route/<resource>/get_<resource>.go`
- `internal/app/api/route/<resource>/get_<resource>s.go`
- `internal/app/api/route/<resource>/update_<resource>.go`
- `internal/app/api/route/<resource>/delete_<resource>.go`

Yêu cầu:

- `route.go` phải đăng ký group `/<resource>s` và gắn `RequireJwt`.
- `get_<resource>.go` phải parse UUID từ `c.Param("id")`. Nếu Resource có trường slug thì sẽ sử dụng `c.Param("slug")` thay vì `c.Param("id")`.
- `get_<resource>s.go` phải dùng `ShouldBindQuery` để bind paging.
- `fx.go` phải dùng `common.ProvideAsRoute(New<Resource>Route)`.
- Cập nhật `internal/app/api/route/fx.go` để thêm module mới.

### Bước 7 - Kiểm tra build

Sau khi tạo xong toàn bộ file, chạy:

```bash
go build ./...
```

Sửa các lỗi compile liên quan đến phần vừa tạo cho đến khi build thành công.

### Bước 8 - Swagger Documentation

Thêm Swagger/OpenAPI annotations vào tất cả handler method files:

- Mỗi handler method phải có comments Swagger phía trên `func`:

```go
// @Summary      Create <resource>
// @Description  Create a new <resource> in the system
// @Tags         <Resource>
// @Accept       json
// @Produce      json
// @Param        body  body      <resource>request.Create<Resource>Request  true  "<Resource> creation request"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     AccessToken
// @Router       /<resource>s [post]
func (h *<Resource>Handler) create<Resource>(c *gin.Context) { ... }
```

- `@Tags`: Dùng tên resource dạng PascalCase (ví dụ: `Tag`, `Category`)
- `@Param`: Loại `path` cho URL param, `query` cho query param, `body` cho JSON body
- `@Success`: Trả về `response.Response` cho detail, `response.PaginationResponse` cho list
- `@Security`: Thêm `AccessToken` cho endpoints có auth. KHÔNG dùng cho public endpoints
- `@Router`: Đường dẫn chính xác từ `route.go`, method lowercase ([get], [post], [put], [delete])

Sau khi thêm annotations, chạy:
```bash
make swagger
# hoặc
swag init -g cmd/dev/main.go -o docs/ --parseDependency --parseInternal
```

Swagger documentation sẽ được generate tự động vào `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`.

## Checklist

- Tạo migration đúng naming convention.
- Tạo model, request DTO, repository, service, route đầy đủ.
- Đăng ký module mới vào các file `fx.go` tổng.
- Giữ đúng naming convention của package, struct, table và route.
- Chỉ sửa các lỗi compile liên quan trực tiếp tới resource mới.
- Thêm Swagger annotations cho tất cả handler methods.
- Chạy `make swagger` để generate API documentation.
