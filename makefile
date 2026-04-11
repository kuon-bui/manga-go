start:
	go run cmd/dev/main.go
migration:
	./scripts/make-migration.sh
migrate:
	sql-migrate up -config=configs/dbconfig.yml
migrate-rollback:
	sql-migrate down -config=configs/dbconfig.yml
migrate-rollback-all:
	sql-migrate down -config=configs/dbconfig.yml -limit=0
swagger:
	swag init -g cmd/dev/main.go -o docs/ --parseDependency --parseInternal
seed:
	go run cmd/seed/main.go

