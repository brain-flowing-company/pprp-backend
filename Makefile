.PHONY: docs migrate

up:
	docker-compose -f docker-compose.dev.yaml up -d --build --no-deps

down:
	docker-compose down

restart:
	docker-compose down && docker-compose -f docker-compose.dev.yaml up -d --build --no-deps

docs:
	swag init -g ./cmd/main.go -o ./docs/

test:
	go test ./internal/... -coverprofile=coverage.out

push:
	docker buildx build -t suechaokhai-backend --platform linux/amd64 .
	docker tag suechaokhai-backend 44.221.177.107:50000/se/suechaokhai-backend
	docker push 44.221.177.107:50000/se/suechaokhai-backend

cover:
	go tool cover -html=coverage.out

migrate:
	atlas schema apply \
		--url "postgres://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable" \
		--to "file://./migrations/schema.sql" \
		--dev-url "docker://postgres/15"