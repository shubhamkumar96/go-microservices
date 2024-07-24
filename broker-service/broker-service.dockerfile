# Base Go image used to build our source code
FROM golang:1.22.4-alpine AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api
RUN chmod +x /app/brokerApp

# Build a tiny docker image
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/brokerApp /app
CMD ["/app/brokerApp"]
