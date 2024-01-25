# ==================== builder ====================
FROM golang:1.21.5-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go

# ==================== runner ====================
FROM alpine:latest as runner
WORKDIR /app

COPY --from=builder /app/main ./cmd/main
COPY --from=builder /app/.env .

CMD ["./cmd/main"]