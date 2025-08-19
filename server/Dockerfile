FROM golang:alpine3.22 AS builder

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN go build -o server_app ..

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/server_app server_app

CMD ["./server_app"]
