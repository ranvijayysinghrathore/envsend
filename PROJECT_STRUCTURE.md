# Project Structure

```
envsend/
├── cli/                          # CLI client (Go)
│   ├── cmd/                      # Cobra commands
│   │   ├── root.go              # Root command
│   │   ├── send.go              # Send command
│   │   └── receive.go           # Receive command
│   ├── crypto/                   # Cryptography layer
│   │   ├── encryption.go        # AES-256-GCM
│   │   ├── keyderivation.go     # Argon2id
│   │   ├── hashing.go           # BLAKE3
│   │   ├── asymmetric.go        # X25519
│   │   ├── shamir.go            # Shamir Secret Sharing
│   │   └── memory.go            # Secure memory handling
│   ├── client/                   # API client
│   │   ├── api.go               # HTTP client
│   │   └── ssh.go               # SSH key fetching
│   ├── utils/                    # Utilities
│   │   └── file.go              # File I/O
│   └── main.go                   # Entry point
│
├── backend/                      # Backend services (Go)
│   ├── cmd/
│   │   ├── server/              # API server
│   │   │   └── main.go
│   │   └── worker/              # Background worker
│   │       └── main.go
│   └── internal/
│       ├── api/                  # HTTP API
│       │   ├── handlers/        # Request handlers
│       │   │   ├── secrets.go
│       │   │   └── health.go
│       │   ├── middleware/      # Middleware
│       │   │   ├── ratelimit.go
│       │   │   ├── security.go
│       │   │   └── audit.go
│       │   └── router.go        # Chi router
│       ├── config/              # Configuration
│       │   └── config.go
│       ├── models/              # Data models
│       │   ├── secret.go
│       │   └── audit.go
│       ├── services/            # Business logic
│       │   └── secret_service.go
│       └── storage/             # Storage layers
│           ├── postgres.go      # PostgreSQL
│           ├── redis.go         # Redis
│           └── s3.go            # MinIO/S3
│
├── migrations/                   # Database migrations
│   ├── 001_initial_schema.up.sql
│   ├── 001_initial_schema.down.sql
│   ├── 002_enterprise_features.up.sql
│   └── 002_enterprise_features.down.sql
│
├── deployments/                  # Deployment configs
│   ├── kubernetes/              # Kubernetes manifests
│   │   ├── namespace.yaml
│   │   ├── configmap.yaml
│   │   ├── api-deployment.yaml
│   │   ├── worker-deployment.yaml
│   │   └── ingress.yaml
│   └── helm/                    # Helm chart
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md          # System architecture
│   ├── CRYPTOGRAPHY.md          # Crypto specification
│   ├── THREAT_MODEL.md          # Security analysis
│   ├── API.md                   # API reference
│   └── DEPLOYMENT.md            # Deployment guide
│
├── Dockerfile.api               # API server Docker build
├── Dockerfile.worker            # Worker Docker build
├── Dockerfile.cli               # CLI Docker build
├── docker-compose.yml           # Local development
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
├── .env.example                 # Environment template
├── .gitignore                   # Git ignore rules
└── README.md                    # Main documentation
```

## File Count: 60+ files

## Lines of Code: ~8,000+ LOC

## Key Technologies

- **Language**: Go 1.21+
- **CLI Framework**: Cobra
- **HTTP Router**: Chi
- **Database**: PostgreSQL
- **Cache**: Redis
- **Storage**: MinIO/S3
- **Container**: Docker
- **Orchestration**: Kubernetes
- **Package Manager**: Helm
