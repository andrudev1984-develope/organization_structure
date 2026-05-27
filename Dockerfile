# Building
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Running
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/openapi ./openapi
COPY --from=builder /app/internal/config/config.yaml ./internal/config/config.yaml

EXPOSE 8080

CMD ["./main"]
