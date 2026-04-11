FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /wishlist-service ./cmd/server/

FROM alpine:latest
RUN apk --no-cache add ca-certificates && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup
WORKDIR /root/
COPY --from=builder /wishlist-service .
COPY --from=builder /app/migrations ./migrations/
RUN chown -R appuser:appgroup /root/
USER appuser
EXPOSE 8080
CMD ["./wishlist-service"]
