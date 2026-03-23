---
name: add-migration
description: Tạo một file migration SQL mới cho dự án manga-go. Bao gồm các bước tạo file, nội dung template và kiểm tra status.

---
---

Tạo một file migration SQL mới cho dự án manga-go.

**Mô tả migration:** ${input:migrationDesc:Mô tả ngắn bằng snake_case. Ví dụ: add_cover_to_mangas, create_tags_table}

---

## Quy trình

1. Lấy timestamp hiện tại bằng lệnh terminal:
   ```bash
   date +%Y%m%d_%H%M%S
   ```

2. Tạo file bằng lệnh: 
   ```bash
   make migration name=<migrationDesc>
   ```

3. Nếu tạo bảng mới, luôn bổ sung:
   - Cột `id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY`
   - Cột `created_at`, `updated_at` kiểu `TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`
   - Cột `deleted_at TIMESTAMPTZ NULL` (soft delete)
   - Trigger `update_<table>_updated_at` gọi `update_updated_at_column()`

4. Phần `-- +migrate Down` phải rollback hoàn toàn phần Up (DROP TRIGGER trước, rồi DROP TABLE).
