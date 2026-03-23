Tạo một file migration SQL mới cho dự án manga-go.

**Mô tả migration:** $ARGUMENTS

---

## Quy trình

1. Lấy timestamp hiện tại:
   ```bash
   date +%Y%m%d_%H%M%S
   ```

2. Tạo file: `migrations/<timestamp>_<mô_tả_snake_case>.sql`

3. Nội dung template:
   ```sql
   -- +migrate Up
   -- SQL thực hiện migration tại đây

   -- +migrate Down
   -- SQL rollback tại đây
   ```

4. Nếu tạo bảng mới, luôn bổ sung:
   - Cột `id uuid NOT NULL DEFAULT uuid_generate_v4()`
   - Cột `created_at`, `updated_at` kiểu `TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`
   - Cột `deleted_at TIMESTAMPTZ NULL` (soft delete)
   - Trigger `update_<table>_updated_at` gọi `update_updated_at_column()`

5. Phần `-- +migrate Down` phải rollback hoàn toàn phần Up (DROP TRIGGER trước, rồi DROP TABLE).

6. Sau khi tạo file, kiểm tra status:
   ```bash
   sql-migrate status -config=configs/dbconfig.yml
   ```
