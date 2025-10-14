FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/api/main.go -o ./docs

RUN go mod tidy && go build -o bookstore-api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bookstore-api .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./bookstore-api"]