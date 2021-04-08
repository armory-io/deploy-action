FROM golang:1.14.2 AS builder

COPY . .
RUN GOPATH= CGO_ENABLED=0 go build -o /bin/action

FROM alpine:3.13.4
RUN apk add ca-certificates bash
COPY --from=builder /bin/action /usr/bin/action
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]