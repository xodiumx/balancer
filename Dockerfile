FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN ls -la .

RUN go build -o balancer ./src/cmd


# Используем минимальный образ для финального контейнера
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/balancer .

EXPOSE 50051
CMD ["./balancer"]
