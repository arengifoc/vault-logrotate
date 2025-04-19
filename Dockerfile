FROM golang:1.18 AS builder
WORKDIR /go/src/github.com/HanseMerkur/vault-logrotate
COPY * ./
RUN go get -d -v \
    && go build .


FROM alpine:3

ENV CRONTAB="0 * * * *"

# Same group/user ids as vault container
RUN apk update --no-cache && apk add --no-cache logrotate && \
    addgroup -g 1000 crond && \
    adduser -u 100 -S -g crond -D -H -h "/tmp" crond

COPY --from=builder /go/src/github.com/HanseMerkur/vault-logrotate/vault-logrotate /usr/local/bin/vault-logrotate

USER crond

ENTRYPOINT ["/usr/local/bin/vault-logrotate"]
