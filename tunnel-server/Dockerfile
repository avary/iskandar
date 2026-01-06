# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY shared/ ./shared/

COPY tunnel-server/ ./tunnel-server/

WORKDIR /build/tunnel-server

# Add shared module as requirement and replace with local path
RUN go mod edit -require=github.com/igneel64/iskandar/shared@v0.0.0
RUN go mod edit -replace=github.com/igneel64/iskandar/shared=../shared
RUN GOWORK=off go mod download
RUN GOWORK=off CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tunnel-server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/tunnel-server .

CMD ["./tunnel-server"]
