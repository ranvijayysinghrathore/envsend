# EnvSend Threat Model

## Assets

### Critical Assets
1. **Plaintext Secrets** - `.env` files, API keys, credentials
2. **Encryption Keys** - AES-256 keys, passphrases
3. **User Privacy** - IP addresses, access patterns
4. **System Availability** - Service uptime, data integrity

### Secondary Assets
1. **Encrypted Blobs** - Ciphertext in S3
2. **Metadata** - Expiry times, view counts
3. **Audit Logs** - Access history

## Threat Actors

### External Attackers
- **Skill Level**: Low to Advanced
- **Motivation**: Data theft, disruption
- **Access**: Internet-facing API

### Malicious Insiders
- **Skill Level**: Advanced
- **Motivation**: Data exfiltration
- **Access**: Server infrastructure

### Compromised Clients
- **Skill Level**: Varies
- **Motivation**: Varies
- **Access**: Client machines

## Attack Vectors

### 1. Server Compromise

**Threat**: Attacker gains access to backend servers

**Impact**: HIGH (but mitigated by zero-knowledge)

**Mitigation**:
- ✅ Zero-knowledge architecture (no plaintext on server)
- ✅ Encrypted blobs only
- ✅ Keys never stored server-side
- ✅ Minimal attack surface
- ✅ Security hardening (SELinux, AppArmor)

**Residual Risk**: LOW - Attacker gets encrypted blobs only

---

### 2. Database Breach

**Threat**: SQL injection or database compromise

**Impact**: MEDIUM

**Mitigation**:
- ✅ Parameterized queries (no SQL injection)
- ✅ Encrypted blobs stored separately (S3)
- ✅ Metadata only in database
- ✅ Audit logs don't contain secrets
- ✅ Database encryption at rest

**Residual Risk**: LOW - Metadata exposure only

---

### 3. Man-in-the-Middle (MITM)

**Threat**: Network interception

**Impact**: HIGH (if successful)

**Mitigation**:
- ✅ TLS 1.3 mandatory
- ✅ HSTS headers
- ✅ Certificate pinning (CLI)
- ✅ Client-side encryption (defense in depth)

**Residual Risk**: LOW - Requires TLS compromise

---

### 4. Replay Attacks

**Threat**: Reuse of intercepted requests

**Impact**: MEDIUM

**Mitigation**:
- ✅ One-time access (max views)
- ✅ Expiry times
- ✅ View counter increment
- ✅ Automatic destruction

**Residual Risk**: LOW - Limited replay window

---

### 5. Brute Force Attacks

**Threat**: Guess passphrases or keys

**Impact**: HIGH (if weak passphrase)

**Mitigation**:
- ✅ Argon2id (memory-hard, slow)
- ✅ 256-bit random keys (2^256 combinations)
- ✅ Rate limiting (100 req/min per IP)
- ✅ No key hints or validation

**Residual Risk**: MEDIUM - Depends on passphrase strength

---

### 6. Phishing

**Threat**: Social engineering to obtain secret URLs

**Impact**: HIGH

**Mitigation**:
- ⚠️ User education
- ⚠️ Short expiry times (default 10 min)
- ⚠️ One-time access
- ❌ Cannot prevent user error

**Residual Risk**: HIGH - Human factor

---

### 7. Client Compromise

**Threat**: Malware on sender/recipient machine

**Impact**: CRITICAL

**Mitigation**:
- ⚠️ Memory zeroing (reduce exposure window)
- ⚠️ No persistent storage
- ❌ Cannot protect compromised clients

**Residual Risk**: CRITICAL - Out of scope

---

### 8. Denial of Service (DoS)

**Threat**: Overwhelm server resources

**Impact**: MEDIUM (availability only)

**Mitigation**:
- ✅ Rate limiting per IP
- ✅ Request size limits (10 MB)
- ✅ Horizontal scaling (Kubernetes HPA)
- ✅ CDN for static content
- ✅ DDoS protection (Cloudflare, etc.)

**Residual Risk**: LOW - Service degradation only

---

### 9. Side-Channel Attacks

**Threat**: Timing attacks, cache attacks

**Impact**: LOW

**Mitigation**:
- ✅ Constant-time crypto (X25519, AES-GCM)
- ✅ No early returns in validation
- ⚠️ Memory zeroing (cache pollution)

**Residual Risk**: LOW - Requires local access

---

### 10. Supply Chain Attacks

**Threat**: Compromised dependencies

**Impact**: CRITICAL

**Mitigation**:
- ✅ Minimal dependencies
- ✅ Vetted crypto libraries (Go stdlib, x/crypto)
- ✅ Dependency pinning (go.sum)
- ✅ Regular security audits
- ⚠️ Automated vulnerability scanning

**Residual Risk**: MEDIUM - Ongoing vigilance required

## Attack Tree

```
[Obtain Plaintext Secret]
├── [Compromise Server]
│   ├── [Exploit API vulnerability] → Mitigated (input validation)
│   ├── [SQL Injection] → Mitigated (parameterized queries)
│   └── [Container escape] → Mitigated (Kubernetes security)
│
├── [Intercept Network Traffic]
│   ├── [MITM attack] → Mitigated (TLS 1.3)
│   └── [DNS hijacking] → Mitigated (DNSSEC, HSTS)
│
├── [Brute Force Decryption]
│   ├── [Guess passphrase] → Mitigated (Argon2id, rate limiting)
│   └── [Guess random key] → Infeasible (2^256 combinations)
│
├── [Social Engineering]
│   ├── [Phish secret URL] → User education
│   └── [Insider threat] → Audit logs, access controls
│
└── [Compromise Client]
    ├── [Malware on sender] → Out of scope
    └── [Malware on recipient] → Out of scope
```

## Security Controls

### Preventive Controls
- Client-side encryption (AES-256-GCM)
- TLS 1.3 (transport security)
- Rate limiting (brute force prevention)
- Input validation (injection prevention)
- Security headers (XSS, clickjacking)

### Detective Controls
- Audit logging (all operations)
- Failed access tracking
- Anomaly detection (future)

### Corrective Controls
- Automatic secret destruction
- Incident response plan
- Security patching process

## Compliance & Standards

### Security Frameworks
- ✅ OWASP Top 10 (2021)
- ✅ NIST Cybersecurity Framework
- ✅ CIS Controls

### Data Protection
- ✅ GDPR-compliant (minimal data collection)
- ✅ CCPA-compliant (no sale of data)
- ⚠️ SOC 2 Type II (future audit)

## Incident Response

### Severity Levels

**Critical**: Plaintext secret exposure
- Response time: Immediate
- Actions: Revoke secrets, notify users, investigate

**High**: Server compromise
- Response time: < 1 hour
- Actions: Isolate systems, forensics, restore

**Medium**: DoS attack
- Response time: < 4 hours
- Actions: Enable DDoS protection, scale resources

**Low**: Failed login attempts
- Response time: < 24 hours
- Actions: Monitor, block IPs if needed

## Security Recommendations

### For Users
1. Use strong passphrases (20+ characters)
2. Verify recipient identity before sharing
3. Use shortest expiry time needed
4. Enable IP locking when possible
5. Don't share secret URLs via insecure channels

### For Operators
1. Enable TLS 1.3 only (disable older versions)
2. Use WAF (Web Application Firewall)
3. Enable DDoS protection
4. Regular security audits
5. Automated vulnerability scanning
6. Incident response plan
7. Regular backups (metadata only)

## Known Limitations

1. **Cannot prevent phishing** - User education required
2. **Cannot protect compromised clients** - Out of scope
3. **Passphrase strength** - User responsibility
4. **URL sharing security** - User responsibility
5. **No multi-factor authentication** - Future enhancement

## Future Security Enhancements

- [ ] Post-quantum cryptography (Kyber, Dilithium)
- [ ] Hardware security module (HSM) support
- [ ] Anomaly detection (ML-based)
- [ ] Canary tokens (honeypots)
- [ ] Security audit by third party
- [ ] Bug bounty program
- [ ] Penetration testing
