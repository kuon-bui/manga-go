---
agent: agent
description: Tạo đầy đủ một CRUD resource mới theo đúng kiến trúc của dự án manga-go.
tools:
  - search/codebase
  - execute/runInTerminal
  - edit/createFile
  - edit/editFiles
---

Tạo đầy đủ CRUD resource mới cho dự án manga-go.

**Resource cần tạo:** ${input:resourceName:Tên resource (singular, PascalCase). Ví dụ: Tag, Category, Publisher}

---

## Quy trình thực hiện

Đọc `CLAUDE.md` để nắm convention, sau đó tạo **toàn bộ** các file sau theo thứ tự:

### Bước 1 — Migration

Tên file: `migrations/<YYYYMMDD>_<HHMMSS>_create_<resource>_table.sql`
Lấy timestamp thực tế từ lệnh terminal: `date +%Y%m%d_%H%M%S`

Template:
```sql
-- +migrate Up
CREATE TABLE <resource>s (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    -- thêm các cột domain tại đây
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

### Bước 2 — Model

File: `internal/pkg/model/<resource>.go`
- Embed `common.SqlModel`
- Implement `TableName() string`

### Bước 3 — Request DTOs

Files:
- `internal/pkg/request/<resource>/create.go` — `Create<Resource>Request`
- `internal/pkg/request/<resource>/update.go` — `Update<Resource>Request`

Dùng tag `binding:"required"` cho các trường bắt buộc.

### Bước 4 — Repository

Files:
- `internal/pkg/repo/<resource>/repo.go` — embed `*base.BaseRepository[model.<Resource>]`
- `internal/pkg/repo/<resource>/fx.go` — `var Module = fx.Module(...)`
- Cập nhật `internal/pkg/repo/fx.go` — thêm module mới vào danh sách

### Bước 5 — Service

Files:
- `internal/pkg/services/<resource>/service.go` — struct + params + constructor
- `internal/pkg/services/<resource>/fx.go` — `var Module = fx.Module(...)`
- `internal/pkg/services/<resource>/create.go`
- `internal/pkg/services/<resource>/get.go` — dùng `errors.Is(err, gorm.ErrRecordNotFound)` để trả `ResultNotFound`
- `internal/pkg/services/<resource>/list.go` — dùng `FindPaginated` + `ResponsePaginationData`
- `internal/pkg/services/<resource>/list_all.go` — dùng `FindAll`
- `internal/pkg/services/<resource>/update.go`
- `internal/pkg/services/<resource>/delete.go`
- Cập nhật `internal/pkg/services/fx.go` — thêm module mới

### Bước 6 — API Route

Files:
- `internal/app/api/route/<resource>/handler.go`
- `internal/app/api/route/<resource>/route.go` — Setup() với group `/<resource>s` + `RequireJwt`
- `internal/app/api/route/<resource>/fx.go`
- `internal/app/api/route/<resource>/create_<resource>.go`
- `internal/app/api/route/<resource>/get_<resource>.go` — parse UUID từ `c.Param("id")`
- `internal/app/api/route/<resource>/get_<resource>s.go` — ShouldBindQuery paging
- `internal/app/api/route/<resource>/get_all_<resource>s.go`
- `internal/app/api/route/<resource>/update_<resource>.go`
- `internal/app/api/route/<resource>/delete_<resource>.go`
- Cập nhật `internal/app/api/route/fx.go` — thêm module mới

### Bước 7 — Kiểm tra

Sau khi tạo xong tất cả file, chạy:
```bash
go build ./...
```

Sửa mọi lỗi compile cho đến khi build thành công.
