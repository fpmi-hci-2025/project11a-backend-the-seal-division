FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go mod tidy && go build -o /bookstore-api ./cmd/api

EXPOSE 8080

CMD ["/bookstore-api"]