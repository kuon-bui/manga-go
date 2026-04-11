.PHONY: \
	start migration migrate migrate-rollback migrate-rollback-all swagger \
	test test-unit test-race test-integration test-integration-verbose lint coverage benchmark clean-test-cache

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
seed-data:
	go run cmd/seed/main.go

test: test-unit

test-unit:
	go test ./... -count=1

test-race:
	go test ./... -race -count=1

test-integration:
	go test ./... -tags=integration -count=1

test-integration-verbose:
	go test ./... -tags=integration -count=1 -v

lint:
	golangci-lint run ./...

coverage:
	go test ./... -covermode=atomic -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out

benchmark:
	go test ./... -run='^$' -bench=. -benchmem -count=1

clean-test-cache:
	go clean -testcache

