# Iskandar

<p align="center">
  <img src="assets/iskndr.png" alt="Iskandar Logo" width="200"/>
</p>

A lightweight, self-hosted HTTP tunnel service for exposing local applications to the internet in Go.

## Features

- 🚀 **Simple CLI** - Expose local apps with a single command
- 📦 **Self-Hosted** - Full control over your infrastructure
- 🔒 **HTTPS Support** - Built-in TLS termination with nginx
- 🌐 **Wildcard Subdomains** - Automatic subdomain allocation for each tunnel

## Quick Start

### Using the CLI

```bash
# Clone the repository
git clone https://github.com/igneel64/iskandar.git
cd iskandar/iskndr

# Build the CLI
go build -o iskndr ./cmd/iskndr

# Expose a local service
./iskndr tunnel --server https://myiskandar.server.deployment.com 3000
```

Replace `https://myiskandar.server.deployment.com` with your tunnel server URL.

## Self-Hosting

Complete deployment instructions with Docker, nginx, and HTTPS setup are available in [tunnel-server/DEPLOYMENT.md](tunnel-server/DEPLOYMENT.md).

## Project Structure

```

iskandar/
├── iskndr/ # CLI client
├── tunnel-server/ # HTTP tunnel server
└── shared/ # Shared code

```

## License

MIT License - see [LICENSE](LICENSE) for details.
