FROM golang:alpine AS builder
ENV CGO_ENABLED=0
ENV CGO_CHECK=0
RUN apk update && \
  apk add --no-cache git=2.36.2-r0
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
RUN go build -o /usr/local/bin/bluey main.go


FROM alpine:3.16
LABEL maintainer="myoung34@my.apsu.edu"

# hadolint ignore=DL3018
RUN apk add -U --no-cache bluez

COPY --from=builder /usr/local/bin/bluey /usr/local/bin/bluey
VOLUME "/etc/bluey"
ENTRYPOINT ["/usr/local/bin/bluey"]
CMD ["-c", "/etc/bluey/config.toml"]
