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

REGISTRY=registry.digitalocean.com/brainflowingcompany
push:
	docker buildx build -t ${REGISTRY}/se/suechaokhai-backend --platform linux/amd64 .
	docker push ${REGISTRY}/se/suechaokhai-backend
	docker image rm ${REGISTRY}/se/suechaokhai-backend

cover:
	go tool cover -html=coverage.out

migrate:
	atlas schema apply \
		--url "postgres://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable" \
		--to "file://./migrations/schema.sql" \
		--dev-url "docker://postgres/15"