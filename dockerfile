FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.2.1@sha256:8879a398dedf0aadaacfbd332b29ff2f84bc39ae6d4e9c0a1109db27ac5ba012 AS xx
FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.20 AS build-stage

## Copy from the xx platform build helper
COPY --from=xx / /

RUN apk add --update --no-cache ca-certificates make git curl clang lld

ARG TARGETPLATFORM

RUN xx-apk --update --no-cache add musl-dev gcc

RUN xx-go --wrap

WORKDIR /usr/local/src/gohlay

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# See https://github.com/confluentinc/confluent-kafka-go#librdkafka
# See https://github.com/confluentinc/confluent-kafka-go#static-builds-on-linux
RUN go build -ldflags '-linkmode external -extldflags "-static"' -tags musl -o /gohlay .
RUN xx-verify /gohlay

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /gohlay /gohlay

ENTRYPOINT ["/gohlay"]
