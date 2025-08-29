FROM golang:alpine3.22 AS builder

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN go build -o server_app .

FROM chromedp/headless-shell:stable

WORKDIR /app

COPY --from=builder /go/src/app/server_app server_app

ENTRYPOINT ["./server_app"]