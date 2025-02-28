FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server ./cmd/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/api_server .
COPY configs/app.toml configs/app.toml

EXPOSE 8080