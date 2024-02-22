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

ARG GOOGLE_CLIENT_SECRET
ENV GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET}

ARG AWS_SECRET_ACCESS_KEY
ENV AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}

ARG EMAIL_PASSWORD
ENV EMAIL_PASSWORD=${EMAIL_PASSWORD}

CMD ["./cmd/main"]