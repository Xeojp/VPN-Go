FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o vpn-service ./cmd/server

FROM alpine:latest
RUN apk add --no-cache wireguard-tools iptables iproute2 sqlite bash
COPY --from=builder /app/vpn-service /vpn-service
COPY web/static /web/static
EXPOSE 8080 51820/udp
CMD ["/vpn-service"]
