# Contributing to Iskandar

Thank you for your interest in contributing to Iskandar!

## Development Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/igneel64/iskandar.git
   cd iskandar
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Run the tunnel server locally**

   ```bash
   cd tunnel-server
   ISKNDR_BASE_DOMAIN=localhost.direct \
   ISKNDR_BASE_SCHEME=http \
   ISKNDR_PORT=8080 \
   ISKNDR_LOGGING=true \
   go run .
   ```

4. **Run the CLI (in another terminal)**
   ```bash
   cd iskndr
   go run ./cmd/iskndr tunnel localhost:3000 --server http://localhost:8080
   ```

### E2E Tests

```bash
cd e2e
go test -tags=e2e -v ./...
```

## Release Process

Releases are automated using GoReleaser and GitHub Actions.

### Creating a Release

1. **Ensure all tests pass** on the main branch

2. **Create and push a version tag**

   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
