FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server ./cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN cp /go/bin/goose /app/goose

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api_server .
COPY --from=builder /app/goose /usr/local/bin/goose
COPY configs/app.toml configs/app.toml
COPY migrations /migrations

EXPOSE 8080