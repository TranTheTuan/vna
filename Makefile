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
	export $(shell cat .env) && \
		migrate -path internal/migrations -database postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable up

migrate-down:
	export $(shell cat .env) && \
		migrate -path internal/migrations -database postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable down $(NUM)

migrate-force:
	export $(shell cat .env) && \
		migrate -path internal/migrations -database postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable force $(VER)
