# Deployment Guide

## Self-Hosting with Docker and Nginx

This guide explains how to deploy the iskandar tunnel server with nginx as a reverse proxy.

### Prerequisites

- Docker and Docker Compose installed
- A domain name (e.g., `tunnel.example.com`)
- DNS configured with wildcard A record: `*.tunnel.example.com` pointing to your server
- SSL certificates (from Let's Encrypt or your provider)
- Ports 80 and 443 open on your firewall

### Docker Compose Setup

Create a `docker-compose.yml` file:

```yaml
services:
  tunnel-server:
    build: .
    container_name: iskandar-tunnel
    restart: unless-stopped
    environment:
      - ISKNDR_BASE_SCHEME=https
      - ISKNDR_BASE_DOMAIN=tunnel.example.com
      - ISKNDR_PORT=8080
    networks:
      - tunnel-net
    expose:
      - "8080"

  nginx:
    image: nginx:alpine
    container_name: iskandar-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./certs:/etc/nginx/certs:ro
    networks:
      - tunnel-net
    depends_on:
      - tunnel-server

networks:
  tunnel-net:
    driver: bridge
```

### Environment Variables

| Variable             | Description                | Default                 |
| -------------------- | -------------------------- | ----------------------- |
| `ISKNDR_BASE_SCHEME` | URL scheme for tunnel URLs | `http`                  |
| `ISKNDR_BASE_DOMAIN` | Base domain for tunnels    | `localhost.direct:8080` |
| `ISKNDR_PORT`        | Port the server listens on | `8080`                  |

### Start the Server

```bash
docker-compose up -d
```

### Testing

Test the tunnel connection:

```bash
./iskndr tunnel --server https://tunnel.example.com 8080
```
