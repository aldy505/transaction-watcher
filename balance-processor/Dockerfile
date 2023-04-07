FROM golang:1.20.3-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build -o balance-processor .

FROM debian:bullseye AS runtime
WORKDIR /app
RUN apt-get update && apt-get install -y curl
COPY --from=builder /app/balance-processor balance-processor

CMD ["/app/balance-processor"]