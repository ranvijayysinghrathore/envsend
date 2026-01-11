# EnvSend API Reference

## Base URL

```
Production: https://envsend.yourdomain.com
Local Dev:  http://localhost:8080
```

## Authentication

No authentication required for public API. Enterprise features require API keys (future).

## Rate Limiting

- **Limit**: 100 requests per minute per IP address
- **Window**: 60 seconds
- **Headers**:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Requests remaining in current window
  - `Retry-After`: Seconds until rate limit resets (when exceeded)

## Endpoints

### Create Secret

**POST** `/api/v1/secrets`

Create a new encrypted secret.

**Request Body**:
```json
{
  "encrypted_blob": "base64-encoded-ciphertext",
  "encryption_metadata": {
    "algorithm": "AES-256-GCM",
    "iv": "base64-encoded-iv",
    "keyDerivation": "none|argon2id|x25519",
    "salt": "base64-encoded-salt",
    "version": "1.0"
  },
  "expires_in": "10m",
  "max_views": 1,
  "ip_lock": "192.168.1.1",
  "recipient_id": "github:username"
}
```

**Request Fields**:
- `encrypted_blob` (required): Base64-encoded encrypted data
- `encryption_metadata` (required): Encryption parameters
- `expires_in` (required): Duration string (e.g., "10m", "1h", "24h")
- `max_views` (required): Maximum number of views (min: 1)
- `ip_lock` (optional): Lock secret to specific IP address
- `recipient_id` (optional): Recipient identifier for SSH encryption

**Response** (201 Created):
```json
{
  "secret_id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "/s/550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "2024-01-15T10:30:00Z",
  "max_views": 1
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body or parameters
- `413 Payload Too Large`: Encrypted blob exceeds 10 MB
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

**Example**:
```bash
curl -X POST https://envsend.io/api/v1/secrets \
  -H "Content-Type: application/json" \
  -d '{
    "encrypted_blob": "SGVsbG8gV29ybGQh...",
    "encryption_metadata": {
      "algorithm": "AES-256-GCM",
      "iv": "cmFuZG9tSVY=",
      "keyDerivation": "none",
      "version": "1.0"
    },
    "expires_in": "10m",
    "max_views": 1
  }'
```

---

### Get Secret

**GET** `/api/v1/secrets/{id}`

Retrieve an encrypted secret.

**Path Parameters**:
- `id`: Secret UUID

**Response** (200 OK):
```json
{
  "encrypted_blob": "base64-encoded-ciphertext",
  "encryption_metadata": {
    "algorithm": "AES-256-GCM",
    "iv": "base64-encoded-iv",
    "keyDerivation": "none",
    "version": "1.0"
  },
  "views_remaining": 0,
  "expires_at": "2024-01-15T10:30:00Z"
}
```

**Response Fields**:
- `encrypted_blob`: Base64-encoded encrypted data
- `encryption_metadata`: Decryption parameters
- `views_remaining`: Number of views left before destruction
- `expires_at`: Expiration timestamp (ISO 8601)

**Error Responses**:
- `404 Not Found`: Secret doesn't exist or has expired
- `410 Gone`: Secret has been destroyed (max views reached)
- `403 Forbidden`: IP address mismatch (if IP lock enabled)
- `429 Too Many Requests`: Rate limit exceeded

**Example**:
```bash
curl https://envsend.io/api/v1/secrets/550e8400-e29b-41d4-a716-446655440000
```

---

### Health Check

**GET** `/health`

Check if the server is healthy.

**Response** (200 OK):
```json
{
  "status": "healthy"
}
```

---

### Readiness Check

**GET** `/ready`

Check if the server is ready to accept traffic.

**Response** (200 OK):
```json
{
  "status": "ready"
}
```

---

### Short URL Redirect

**GET** `/s/{id}`

User-friendly redirect to secret retrieval (for web UI).

**Path Parameters**:
- `id`: Secret UUID

**Response**: 307 Temporary Redirect to `/api/v1/secrets/{id}`

---

## Error Format

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

## Status Codes

- `200 OK`: Request successful
- `201 Created`: Secret created successfully
- `307 Temporary Redirect`: Short URL redirect
- `400 Bad Request`: Invalid request
- `403 Forbidden`: Access denied (IP lock)
- `404 Not Found`: Secret not found or expired
- `410 Gone`: Secret destroyed
- `413 Payload Too Large`: Request too large
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## Security Headers

All responses include:

```
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'none'; frame-ancestors 'none'
Referrer-Policy: no-referrer
```

## CORS

CORS is enabled for all origins by default (configurable).

**Headers**:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## Encryption Metadata Format

### None (Random Key)

```json
{
  "algorithm": "AES-256-GCM",
  "iv": "base64-encoded-12-byte-iv",
  "keyDerivation": "none",
  "version": "1.0"
}
```

### Argon2id (Passphrase)

```json
{
  "algorithm": "AES-256-GCM",
  "iv": "base64-encoded-12-byte-iv",
  "keyDerivation": "argon2id",
  "salt": "base64-encoded-16-byte-salt",
  "version": "1.0"
}
```

### X25519 (SSH Key Exchange)

```json
{
  "algorithm": "AES-256-GCM",
  "iv": "base64-encoded-12-byte-iv",
  "keyDerivation": "x25519",
  "encryptedKey": "base64-encoded-encrypted-symmetric-key",
  "ephemeralPubKey": "base64-encoded-32-byte-public-key",
  "version": "1.0"
}
```

## Duration Format

Durations use Go's time.Duration format:

- `s`: seconds
- `m`: minutes
- `h`: hours

**Examples**:
- `10m`: 10 minutes
- `1h`: 1 hour
- `30s`: 30 seconds
- `24h`: 24 hours

**Maximum**: 168h (7 days)

## Limits

- **Max Secret Size**: 10 MB (encrypted)
- **Max Expiry**: 168 hours (7 days)
- **Max Views**: No limit (but recommended < 100)
- **Rate Limit**: 100 requests/minute per IP

## Best Practices

1. **Always use HTTPS** in production
2. **Set shortest expiry** needed
3. **Use max_views: 1** for one-time secrets
4. **Enable IP lock** when recipient IP is known
5. **Verify encryption_metadata** before decryption
6. **Handle 410 Gone** gracefully (secret destroyed)
7. **Respect rate limits** (check headers)

## Client Libraries

### Official
- Go CLI (this repository)

### Community
- JavaScript/TypeScript (coming soon)
- Python (coming soon)
- Rust (coming soon)

## Changelog

### v1.0.0 (2024-01-15)
- Initial release
- Basic secret create/retrieve
- Rate limiting
- Health checks
