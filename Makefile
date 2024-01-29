.PHONY: docs

up:
	docker-compose -f docker-compose.dev.yaml up -d --build --no-deps

docs:
	swag init -g ./cmd/main.go -o ./docs/

test:
	go test ./internal/... -coverprofile=coverage.out

cover:
	go tool cover -html=coverage.out