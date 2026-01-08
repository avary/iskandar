# Test App Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY test-app/go.mod test-app/go.sum ./

RUN go mod download

# Copy source
COPY test-app/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o test-app .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

COPY --from=builder /build/test-app .
COPY test-app/views ./views
COPY test-app/public ./public

EXPOSE 3003

CMD ["./test-app"]
