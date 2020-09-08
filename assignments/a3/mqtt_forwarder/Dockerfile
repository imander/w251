FROM golang:1.15 AS builder

WORKDIR /go/src/mids.w251.edu/a3/mqtt-forwarder
COPY . .
RUN CGO_ENABLED=0 go build -o /tmp/mqtt-forwarder main.go

FROM alpine

COPY --from=builder /tmp/mqtt-forwarder /bin/mqtt-forwarder
ENTRYPOINT ["/bin/mqtt-forwarder"]
