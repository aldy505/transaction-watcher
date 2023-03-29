FROM rust:1.68.2-bullseye AS builder
WORKDIR /app
COPY . .
RUN cargo build --release

FROM debian:bullseye AS runtime
WORKDIR /app
RUN apt-get update && apt-get install -y curl
COPY --from=builder /app/target/release/inserter inserter

CMD ["/app/inserter"]