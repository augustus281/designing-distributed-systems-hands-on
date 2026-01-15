FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY main.go .

RUN go build -o topz main.go

FROM alpine:latest

COPY --from=builder /app/topz /topz

ENTRYPOINT ["/topz"]
