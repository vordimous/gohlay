FROM golang:alpine

RUN apk update && apk add librdkafka-dev pkgconf

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-linkmode external -extldflags -static" -a main.go

EXPOSE 8000

ENTRYPOINT ["/app/main"]
