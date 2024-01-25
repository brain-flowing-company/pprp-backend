.PHONY: docs

up:
	docker-compose up -d

docs:
	swag init -g ./cmd/main.go -o ./docs/

test:
	go test ./internal/... -coverprofile=coverage.out

cover:
	go tool cover -html=coverage.out