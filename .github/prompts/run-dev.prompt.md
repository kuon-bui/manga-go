---
agent: agent
description: Khởi động môi trường dev cho dự án manga-go.
tools:
  - terminal
---

Khởi động môi trường dev cho dự án manga-go theo các bước sau:

## Bước 1 — Kiểm tra infrastructure

Đảm bảo Docker services đang chạy (PostgreSQL, Redis, Jaeger, Mailpit):
```bash
docker compose ps
```

Nếu chưa chạy, khởi động:
```bash
docker compose up -d
```

## Bước 2 — Kiểm tra migration

```bash
sql-migrate status -config=configs/dbconfig.yml
```

Nếu có migration chưa apply:
```bash
sql-migrate up -config=configs/dbconfig.yml
```

## Bước 3 — Kiểm tra build

```bash
go build ./...
```

Sửa lỗi compile nếu có.

## Bước 4 — Khởi động API server với live reload

```bash
air -c api.air.toml
```

## Bước 5 — (Tùy chọn) Khởi động queue worker

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
