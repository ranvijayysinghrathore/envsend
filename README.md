# EnvSend - Zero-Knowledge Secret Transfer Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

**EnvSend** is a production-grade, CLI-first platform for securely transferring `.env` files and secrets between developers without ever exposing plaintext to the server.

## Core Philosophy

- **Zero-knowledge by design** - Server never sees plaintext secrets
- **Client-side encryption only** - AES-256-GCM with Argon2id
- **Ephemeral secrets by default** - Auto-expire and self-destruct
- **CLI-first UX** - Built for developers, by developers
- **Horizontally scalable** - Stateless architecture
- **Secure by default** - No configuration required

## Quick Start

### Installation

```bash
# Download CLI binary
curl -sSL https://envsend.io/install.sh | bash

# Or build from source
git clone https://github.com/ranvijayysinghrathore/envsend.git
cd envsend
make build-cli
sudo make install-cli
```

### Basic Usage

```bash
# Send a secret (default: 10 minutes, 1 view)
envsend .env

# Send with custom expiry and views
envsend .env --expires 1h --max-views 3

# Passphrase-protected
envsend .env --require-passphrase

# Encrypt for specific recipient (SSH key exchange)
envsend .env --ssh github:username

# Pipe support
cat .env | envsend

# ☢️ Nuclear Mode (Shamir Secret Sharing)
# Split secret into 5 shares, requiring any 3 to decrypt
envsend .env --shamir-shares 5 --shamir-threshold 3

# Receive a secret (Smart Mode - detect URL automatically)
envsend "http://localhost:8080/s/..." > .env

# Legacy receive command
envsend receive "http://localhost:8080/s/..." > .env
envreceive https://envsend.io/s/abc123#key > .env
```

## Architecture

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Sender    │         │   EnvSend    │         │  Recipient  │
│   (CLI)     │────────▶│   Server     │◀────────│   (CLI)     │
└─────────────┘         └──────────────┘         └─────────────┘
      │                        │                        │
      │ 1. Encrypt locally     │                        │
      │ 2. Upload blob         │                        │
      │                        │ 3. Store encrypted     │
      │                        │    (PostgreSQL + S3)   │
      │                        │                        │
      │                        │ 4. Download encrypted  │
      │                        │                        │
      │                        │                  5. Decrypt locally
      │                        │                        │
      └────────────────────────┴────────────────────────┘
           Server NEVER sees plaintext secrets
```

### Components

- **CLI Client** (Go) - Encryption, decryption, key management
- **API Server** (Go + Chi) - Metadata management, blob storage
- **Worker** (Go) - TTL cleanup, expired secret deletion
- **PostgreSQL** - Secret metadata and audit logs
- **Redis** - Rate limiting, distributed locks
- **MinIO/S3** - Encrypted blob storage

## Cryptography

All encryption happens **client-side only**:

- **Symmetric Encryption**: AES-256-GCM
- **Key Derivation**: Argon2id (OWASP recommended)
- **Hashing**: BLAKE3
- **Asymmetric Encryption**: X25519 (for SSH key exchange)
- **Secret Sharing**: Shamir's Secret Sharing (The "Nuclear Code" mode)
- **Memory Security**: Explicit zeroing after use

## Deployment

### Local Development

```bash
# Start infrastructure (PostgreSQL, Redis, MinIO)
docker-compose up -d

# Run migrations
make migrate-up

# Start API server
make run-api

# Start worker
make run-worker
```

### Production (Kubernetes)

```bash
# Deploy with kubectl
kubectl apply -f deployments/kubernetes/

# Or use Helm
helm install envsend deployments/helm/ --namespace envsend
```

### Environment Variables

See [`.env.example`](.env.example) for all configuration options.

## Security Features

- ✅ Zero-knowledge architecture
- ✅ Client-side encryption only
- ✅ No plaintext logging
- ✅ Rate limiting per IP
- ✅ IP address locking (optional)
- ✅ Automatic secret expiry
- ✅ One-time access enforcement
- ✅ Comprehensive audit logs
- ✅ Secure memory handling
- ✅ HTTPS/TLS only in production

## Scalability

- **Stateless API** - Horizontal scaling via Kubernetes HPA
- **Object Storage** - Infinite encrypted blob capacity
- **Background Workers** - Distributed TTL cleanup
- **Redis Caching** - Fast rate limiting and locks
- **CDN-Ready** - Encrypted blobs can be cached

## Development

```bash
# Run tests
make test

# Run linters
make lint

# Format code
make fmt

# Build all binaries
make build
```

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [Cryptography](docs/CRYPTOGRAPHY.md)
- [Threat Model](docs/THREAT_MODEL.md)
- [API Reference](docs/API.md)
- [Deployment Guide](docs/DEPLOYMENT.md)

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) first.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Chi](https://github.com/go-chi/chi) - HTTP router
- [MinIO](https://min.io/) - Object storage
- [BLAKE3](https://github.com/BLAKE3-team/BLAKE3) - Cryptographic hashing

## Security Disclosure

Found a security issue? Please contact us instead of opening a public issue.

---

**Made with ❤️ for developers who care about security**
