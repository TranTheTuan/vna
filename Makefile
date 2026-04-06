.PHONY: run build tidy migrate migrate-down migrate-force docs

run:
	go run ./cmd/api/...

build:
	go build -o bin/api ./cmd/api/...

tidy:
	go mod tidy

docs:
	swag init -g cmd/api/main.go -o internal/docs

migrate:
	export $(shell grep -v '^#' .env | xargs) && \
		migrate -path internal/migrations -database "postgres://$$DATABASE_USER:$$DATABASE_PASSWORD@$$DATABASE_HOST:$$DATABASE_PORT/$$DATABASE_DBNAME?sslmode=disable" up

migrate-dev:
	export $(shell grep -v '^#' .env.dev | xargs) && \
		migrate -path internal/migrations -database "postgres://$$DATABASE_USER:$$DATABASE_PASSWORD@$$DATABASE_HOST:$$DATABASE_PORT/$$DATABASE_DBNAME?sslmode=disable" up

migrate-down:
	export $(shell grep -v '^#' .env | xargs) && \
		migrate -path internal/migrations -database "postgres://$$DATABASE_USER:$$DATABASE_PASSWORD@$$DATABASE_HOST:$$DATABASE_PORT/$$DATABASE_DBNAME?sslmode=disable" down $(NUM)

migrate-force:
	export $(shell grep -v '^#' .env | xargs) && \
		migrate -path internal/migrations -database "postgres://$$DATABASE_USER:$$DATABASE_PASSWORD@$$DATABASE_HOST:$$DATABASE_PORT/$$DATABASE_DBNAME?sslmode=disable" force $(VER)
