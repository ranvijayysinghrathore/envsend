# EnvSend Deployment Guide

## Prerequisites

### Local Development
- Docker & Docker Compose
- Go 1.21+
- Make
- PostgreSQL client (for migrations)

### Production
- Kubernetes cluster (1.25+)
- kubectl configured
- Helm 3.x (optional)
- Domain name with DNS access
- TLS certificates (Let's Encrypt recommended)

## Local Development Setup

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/envsend.git
cd envsend
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL, Redis, and MinIO
docker-compose up -d

# Wait for services to be ready
sleep 10
```

### 3. Run Migrations

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
make migrate-up
```

### 4. Build and Run

```bash
# Build all binaries
make build

# Run API server
./bin/api-server

# In another terminal, run worker
./bin/worker

# In another terminal, test CLI
./bin/envsend --help
```

### 5. Test the System

```bash
# Create a test secret
echo "TEST_KEY=test_value" > test.env
./bin/envsend test.env

# Copy the URL and retrieve it
./bin/envsend receive <url> > retrieved.env

# Verify
cat retrieved.env
```

## Production Deployment (Kubernetes)

### 1. Prepare Secrets

Create a `secrets.yaml` file (DO NOT commit to git):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: envsend-secrets
  namespace: envsend
type: Opaque
stringData:
  database-url: "postgres://user:password@postgres:5432/envsend?sslmode=require"
  redis-url: "redis://redis:6379/0"
  s3-endpoint: "minio:9000"
  s3-access-key: "your-access-key"
  s3-secret-key: "your-secret-key"
```

Apply the secret:
```bash
kubectl apply -f secrets.yaml
```

### 2. Build Docker Images

```bash
# Set your registry
export DOCKER_REGISTRY=your-registry.com

# Build images
make docker-build

# Push to registry
make docker-push
```

### 3. Update Kubernetes Manifests

Edit `deployments/kubernetes/api-deployment.yaml` and `worker-deployment.yaml`:

```yaml
image: your-registry.com/envsend-api:latest
```

### 4. Deploy to Kubernetes

```bash
# Create namespace
kubectl apply -f deployments/kubernetes/namespace.yaml

# Apply ConfigMap
kubectl apply -f deployments/kubernetes/configmap.yaml

# Apply secrets (created in step 1)
kubectl apply -f secrets.yaml

# Deploy PostgreSQL (or use managed service)
kubectl apply -f deployments/kubernetes/postgres-statefulset.yaml

# Deploy Redis
kubectl apply -f deployments/kubernetes/redis-deployment.yaml

# Deploy MinIO (or use S3)
kubectl apply -f deployments/kubernetes/minio-statefulset.yaml

# Deploy API server
kubectl apply -f deployments/kubernetes/api-deployment.yaml

# Deploy worker
kubectl apply -f deployments/kubernetes/worker-deployment.yaml

# Deploy ingress
kubectl apply -f deployments/kubernetes/ingress.yaml
```

### 5. Verify Deployment

```bash
# Check pods
kubectl get pods -n envsend

# Check services
kubectl get svc -n envsend

# Check ingress
kubectl get ingress -n envsend

# View logs
kubectl logs -n envsend -l app=envsend-api
kubectl logs -n envsend -l app=envsend-worker
```

### 6. Configure DNS

Point your domain to the ingress load balancer:

```bash
# Get ingress IP
kubectl get ingress -n envsend envsend-ingress

# Create A record
envsend.yourdomain.com -> <INGRESS_IP>
```

### 7. Enable TLS

Install cert-manager:

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

Create ClusterIssuer:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
```

Apply:
```bash
kubectl apply -f clusterissuer.yaml
```

## Helm Deployment (Alternative)

### 1. Install with Helm

```bash
helm install envsend deployments/helm/ \
  --namespace envsend \
  --create-namespace \
  --set image.registry=your-registry.com \
  --set ingress.host=envsend.yourdomain.com \
  --set postgresql.auth.password=your-password
```

### 2. Upgrade Deployment

```bash
helm upgrade envsend deployments/helm/ \
  --namespace envsend \
  --set image.tag=v1.1.0
```

### 3. Uninstall

```bash
helm uninstall envsend --namespace envsend
```

## Managed Services (Recommended for Production)

### PostgreSQL
- AWS RDS
- Google Cloud SQL
- Azure Database for PostgreSQL
- DigitalOcean Managed Databases

### Redis
- AWS ElastiCache
- Google Cloud Memorystore
- Azure Cache for Redis
- Redis Cloud

### Object Storage
- AWS S3
- Google Cloud Storage
- Azure Blob Storage
- DigitalOcean Spaces

### Example with AWS

```bash
# Set environment variables
export DATABASE_URL="postgres://user:pass@rds-endpoint:5432/envsend"
export REDIS_URL="redis://elasticache-endpoint:6379"
export S3_ENDPOINT="s3.amazonaws.com"
export S3_BUCKET="envsend-secrets-prod"
export S3_ACCESS_KEY="AWS_ACCESS_KEY"
export S3_SECRET_KEY="AWS_SECRET_KEY"
export S3_USE_SSL="true"
```

## Monitoring & Observability

### Prometheus Metrics

Add Prometheus annotations to deployments:

```yaml
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "9090"
  prometheus.io/path: "/metrics"
```

### Logging

Use Fluent Bit or Fluentd to collect logs:

```bash
kubectl apply -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/fluent-bit-service-account.yaml
kubectl apply -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/fluent-bit-role.yaml
kubectl apply -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/fluent-bit-role-binding.yaml
kubectl apply -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/output/elasticsearch/fluent-bit-configmap.yaml
kubectl apply -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/output/elasticsearch/fluent-bit-ds.yaml
```

## Scaling

### Horizontal Pod Autoscaling

Already configured in `api-deployment.yaml`:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: envsend-api-hpa
spec:
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        averageUtilization: 70
```

### Database Scaling

- Use read replicas for audit logs
- Connection pooling (already configured)
- Vertical scaling for write-heavy workloads

### Object Storage

- Infinite horizontal scaling (S3/MinIO)
- Enable CDN for encrypted blobs (optional)

## Backup & Disaster Recovery

### Database Backups

```bash
# Automated backups (PostgreSQL)
kubectl exec -n envsend postgres-0 -- pg_dump -U envsend envsend > backup.sql

# Restore
kubectl exec -i -n envsend postgres-0 -- psql -U envsend envsend < backup.sql
```

### Object Storage Backups

- Enable versioning in S3/MinIO
- Cross-region replication (optional)
- Lifecycle policies for old blobs

## Security Hardening

### Network Policies

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: envsend-api-policy
  namespace: envsend
spec:
  podSelector:
    matchLabels:
      app: envsend-api
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
```

### Pod Security Standards

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: envsend
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

## Troubleshooting

### API Server Not Starting

```bash
# Check logs
kubectl logs -n envsend -l app=envsend-api

# Common issues:
# - Database connection failed
# - Redis connection failed
# - MinIO bucket not created
```

### Worker Not Cleaning Up

```bash
# Check worker logs
kubectl logs -n envsend -l app=envsend-worker

# Verify database connectivity
kubectl exec -n envsend postgres-0 -- psql -U envsend -c "SELECT COUNT(*) FROM secrets WHERE destroyed = true;"
```

### Rate Limiting Issues

```bash
# Check Redis connectivity
kubectl exec -n envsend redis-0 -- redis-cli ping

# Clear rate limit for IP
kubectl exec -n envsend redis-0 -- redis-cli DEL "ratelimit:192.168.1.1"
```

## Performance Tuning

### Database

```sql
-- Optimize queries
CREATE INDEX CONCURRENTLY idx_secrets_expires_at_not_destroyed 
ON secrets(expires_at) WHERE NOT destroyed;

-- Analyze tables
ANALYZE secrets;
ANALYZE audit_logs;
```

### Redis

```bash
# Increase max memory
kubectl edit configmap redis-config -n envsend

# Add:
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### API Server

```yaml
# Increase resources
resources:
  requests:
    memory: "256Mi"
    cpu: "200m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

## Maintenance

### Rolling Updates

```bash
# Update image
kubectl set image deployment/envsend-api -n envsend api=your-registry.com/envsend-api:v1.1.0

# Check rollout status
kubectl rollout status deployment/envsend-api -n envsend

# Rollback if needed
kubectl rollout undo deployment/envsend-api -n envsend
```

### Database Migrations

```bash
# Run new migrations
kubectl exec -n envsend api-pod-name -- migrate -path /migrations -database "$DATABASE_URL" up
```

## Cost Optimization

- Use spot instances for workers
- Enable autoscaling (scale to zero during low traffic)
- Use S3 lifecycle policies (delete old encrypted blobs)
- Use managed services (reduce operational overhead)

## Support

For issues and questions:
- GitHub Issues: https://github.com/yourusername/envsend/issues
- Documentation: https://docs.envsend.io
- Email: support@envsend.io
