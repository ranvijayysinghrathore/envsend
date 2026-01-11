# EnvSend Architecture

## System Overview

EnvSend is a **zero-knowledge secret transfer platform** designed for secure, ephemeral sharing of environment files and credentials between developers.

### Core Principles

1. **Zero-Knowledge**: Server never has access to plaintext secrets or decryption keys
2. **Client-Side Encryption**: All cryptographic operations happen on the client
3. **Ephemeral by Default**: Secrets auto-expire and self-destruct
4. **Stateless Backend**: Horizontally scalable API servers
5. **Immutable Audit Trail**: All operations are logged (without secret content)

## Component Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                           │
├────────────────────────────────────────────────────────────────┤
│  CLI (Go)                                                      │
│  ├── Encryption (AES-256-GCM)                                 │
│  ├── Key Derivation (Argon2id)                               │
│  ├── Hashing (BLAKE3)                                        │
│  ├── Asymmetric Crypto (X25519)                              │
│  └── Secure Memory Management                                │
└────────────────────────────────────────────────────────────────┘
                            ▲ │
                            │ │ HTTPS/TLS
                            │ ▼
┌────────────────────────────────────────────────────────────────┐
│                         API LAYER                              │
├────────────────────────────────────────────────────────────────┤
│  API Server (Go + Chi)                                        │
│  ├── REST Endpoints                                           │
│  ├── Rate Limiting (Redis)                                    │
│  ├── Security Headers                                         │
│  ├── Audit Logging                                            │
│  └── Request Validation                                       │
└────────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┼───────────┐
                │           │           │
                ▼           ▼           ▼
┌──────────────────┐ ┌──────────────┐ ┌──────────────┐
│   PostgreSQL     │ │    Redis     │ │   MinIO/S3   │
├──────────────────┤ ├──────────────┤ ├──────────────┤
│ • Metadata       │ │ • Rate Limit │ │ • Encrypted  │
│ • Audit Logs     │ │ • Locks      │ │   Blobs      │
│ • Expiry Times   │ │ • Cache      │ │ • Immutable  │
└──────────────────┘ └──────────────┘ └──────────────┘
                            ▲
                            │
┌────────────────────────────────────────────────────────────────┐
│                      WORKER LAYER                              │
├────────────────────────────────────────────────────────────────┤
│  Background Worker (Go)                                       │
│  ├── TTL Cleanup                                              │
│  ├── Expired Secret Deletion                                  │
│  └── Blob Garbage Collection                                  │
└────────────────────────────────────────────────────────────────┘
```

## Data Flow

### Secret Creation

1. **Client reads file** → `.env` file or stdin
2. **Generate encryption key** → Random 32-byte key (AES-256)
3. **Encrypt locally** → AES-256-GCM with random IV
4. **Upload encrypted blob** → POST to `/api/v1/secrets`
5. **Server stores**:
   - Encrypted blob → MinIO/S3
   - Metadata → PostgreSQL (expiry, max views, etc.)
   - Audit log → PostgreSQL
6. **Return URL** → `https://envsend.io/s/{id}#{key}`
   - Secret ID in path
   - Encryption key in URL fragment (never sent to server)

### Secret Retrieval

1. **Client parses URL** → Extract secret ID and key
2. **Download encrypted blob** → GET `/api/v1/secrets/{id}`
3. **Server**:
   - Validates expiry and view count
   - Checks IP lock (if set)
   - Increments view counter
   - Returns encrypted blob + metadata
   - Destroys if max views reached
4. **Client decrypts locally** → AES-256-GCM with key from URL
5. **Zero memory** → Explicit cleanup of sensitive data

## Security Architecture

### Threat Model

**What we protect against:**
- ✅ Server compromise (zero-knowledge)
- ✅ Database breach (only encrypted blobs)
- ✅ Network interception (TLS + client-side encryption)
- ✅ Replay attacks (one-time access, expiry)
- ✅ Brute force (rate limiting)
- ✅ Memory dumps (explicit zeroing)

**What we don't protect against:**
- ❌ Compromised client machine
- ❌ Malicious sender/recipient
- ❌ Phishing of secret URLs
- ❌ Side-channel attacks on client

### Encryption Layers

1. **Transport Layer**: TLS 1.3 (HTTPS)
2. **Application Layer**: AES-256-GCM (client-side)
3. **Key Derivation**: Argon2id (if passphrase used)
4. **Integrity**: BLAKE3 hashing

## Scalability Design

### Horizontal Scaling

- **API Servers**: Stateless, scale via Kubernetes HPA
- **Workers**: Multiple instances with distributed locks
- **Database**: Read replicas for audit logs
- **Object Storage**: Infinite capacity (MinIO/S3)
- **Cache**: Redis cluster for rate limiting

### Performance Targets

- **API Latency**: < 100ms (p95)
- **Upload Speed**: Limited by network, not server
- **Download Speed**: CDN-cacheable encrypted blobs
- **Throughput**: 10,000+ secrets/second
- **Cleanup**: 100 secrets/second per worker

## Deployment Topology

### Local Development
```
Docker Compose
├── PostgreSQL (port 5432)
├── Redis (port 6379)
├── MinIO (port 9000)
├── API Server (port 8080)
└── Worker (background)
```

### Production (Kubernetes)
```
Kubernetes Cluster
├── Namespace: envsend
├── API Deployment (3-10 replicas, HPA)
├── Worker Deployment (2 replicas)
├── PostgreSQL StatefulSet
├── Redis Deployment
├── MinIO StatefulSet
├── Ingress (TLS termination)
└── Services (ClusterIP)
```

## Database Schema

### Secrets Table
```sql
- id (UUID, PK)
- encrypted_blob_url (TEXT) → S3 reference
- encryption_metadata (JSONB) → Algorithm, IV, etc.
- expires_at (TIMESTAMP)
- max_views (INT)
- view_count (INT)
- ip_lock (INET, nullable)
- destroyed (BOOLEAN)
- created_at, accessed_at, destroyed_at
```

### Audit Logs Table
```sql
- id (BIGSERIAL, PK)
- secret_id (UUID, FK)
- action (ENUM: created, viewed, destroyed, expired)
- ip_address (INET)
- user_agent (TEXT)
- timestamp (TIMESTAMP)
```

## API Endpoints

### Public API

- `POST /api/v1/secrets` - Create secret
- `GET /api/v1/secrets/{id}` - Retrieve secret
- `GET /health` - Health check
- `GET /ready` - Readiness probe

### Short URLs

- `GET /s/{id}` - User-friendly redirect

## Monitoring & Observability

- **Metrics**: Prometheus-compatible endpoints
- **Logging**: Structured JSON logs
- **Tracing**: OpenTelemetry (optional)
- **Alerts**: Expiry failures, storage errors, rate limit breaches

## Future Enhancements

- [ ] Web UI for browser-based decryption
- [ ] QR code generation for mobile
- [ ] Multi-recipient encryption
- [ ] Organization accounts
- [ ] SSO integration (OAuth2/SAML)
- [ ] Compliance reporting (SOC 2, GDPR)
