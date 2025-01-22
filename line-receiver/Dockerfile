
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY ./config ./config
COPY ./handler ./handler
COPY ./middleware ./middleware
COPY ./main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
CMD ["/app/main"]