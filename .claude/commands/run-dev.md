Khởi động môi trường dev cho dự án manga-go.

---

## Các bước thực hiện

### 1. Kiểm tra infrastructure

Đảm bảo Docker services đang chạy (PostgreSQL, Redis, Jaeger, Mailpit):
```bash
docker compose ps
```

Nếu chưa chạy, khởi động:
```bash
docker compose up -d
```

### 2. Kiểm tra migration

```bash
sql-migrate status -config=configs/dbconfig.yml
```

Nếu có migration chưa apply:
```bash
sql-migrate up -config=configs/dbconfig.yml
```

### 3. Kiểm tra build

```bash
go build ./...
```

Sửa lỗi compile nếu có.

### 4. Khởi động API server với live reload

```bash
air -c api.air.toml
```

### 5. (Tùy chọn) Khởi động queue worker

Nếu cần xử lý async tasks (email, background jobs):
```bash
go run cmd/queue/main.go
```

---

## URLs hữu ích

| Service | URL |
|---|---|
| API server | http://localhost:8080 |
| pgAdmin | http://localhost:8081 |
| Asynqmon (queue monitor) | http://localhost:8091 |
| Mailpit (mail UI) | http://localhost:8025 |
| Jaeger (tracing) | http://localhost:16686 |
