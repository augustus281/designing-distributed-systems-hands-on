FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main main.go

CMD ["./main"]