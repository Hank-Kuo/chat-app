FROM golang:1.21-alpine AS go_kafka_base

RUN apk update && apk add --no-cache  \
        bash              \
        build-base        \
        coreutils         \
        gcc               \
        git               \
        make              

RUN apk add librdkafka-dev pkgconf

FROM go_kafka_base as build_base
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_base AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic,netgo -o /bin/app -ldflags "-w -s" ./cmd/api/main.go


FROM alpine:3.10
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/app /bin/app
EXPOSE 9000
EXPOSE 9001
ENTRYPOINT ["/bin/app"]
