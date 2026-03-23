---
agent: agent
description: Chạy database migrations cho dự án manga-go.
tools:
  - terminal
---

Chạy database migrations cho dự án manga-go.

**Hành động:** ${input:action:Chọn hành động: up (apply migrations mới), down (rollback), status (xem trạng thái), redo (rollback + re-apply)}

---

## Các lệnh

### Xem trạng thái hiện tại
```bash
sql-migrate status -config=configs/dbconfig.yml
```

### Apply tất cả migrations mới
```bash
sql-migrate up -config=configs/dbconfig.yml
```

### Rollback migration cuối
```bash
sql-migrate down -config=configs/dbconfig.yml
```

### Redo (rollback + re-apply) migration cuối
```bash
sql-migrate redo -config=configs/dbconfig.yml
```

---

## Lưu ý

- Tool dùng: `sql-migrate` (không phải `goose` hay `migrate`)
- Config file: `configs/dbconfig.yml`
- Thư mục migration: `migrations/`
- Format tên file: `YYYYMMDD_HHMMSS_<mô_tả>.sql`
- Mỗi file có section `-- +migrate Up` và `-- +migrate Down`
- Migrations chạy theo thứ tự timestamp tăng dần
- Sau khi apply, kiểm tra lại bằng `sql-migrate status`
