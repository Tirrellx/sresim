# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o sresim ./cmd

# Run stage
FROM alpine:latest
COPY --from=builder /app/sresim /sresim
EXPOSE 8080
ENTRYPOINT ["/sresim"]
