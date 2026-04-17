# Operations Runbook

> Revision: 2026-04-17 | Owner: Platform Engineering

---

## Table of Contents

1. [On-Call Contacts](#on-call-contacts)
2. [Architecture Overview](#architecture-overview)
3. [Common Alerts & Remediation](#common-alerts--remediation)
4. [Deployment & Rollback](#deployment--rollback)
5. [Disaster Recovery](#disaster-recovery)
6. [Database Procedures](#database-procedures)
7. [Secret Rotation](#secret-rotation)
8. [Audit Log Verification](#audit-log-verification)

---

## On-Call Contacts

| Role | Contact | Escalation |
|------|---------|------------|
| Primary On-Call | Refer to PagerDuty rotation | — |
| Security incidents | security@example.com | CISO |
| Database | DBA team Slack `#dba-oncall` | Head of Infrastructure |

---

## Architecture Overview

```
Internet → Load Balancer (TLS termination)
         → gin-vue-admin (Go/Gin, port 8888)
              ├── PostgreSQL / MySQL (GORM)
              ├── Redis (JWT blacklist, distributed locks, session)
              └── OTEL Collector → Jaeger / Grafana Tempo
```

Key environment variables:

| Variable | Purpose |
|----------|---------|
| `GVA_JWT_SIGNING_KEY` | JWT HMAC secret (≥ 32 chars, never default) |
| `GVA_DB_PASSWORD` | Database password override |
| `GVA_REDIS_PASSWORD` | Redis password override |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector address (default `localhost:4317`) |
| `OTEL_SERVICE_NAME` | Service name in traces (default `gin-vue-admin`) |
| `GIN_MODE` | `release` for production |

---

## Common Alerts & Remediation

### `HighErrorRate` (5xx > 0.5%)

1. Check application logs: `kubectl logs -l app=gin-vue-admin --since=10m`
2. Check DB connectivity: `kubectl exec -it <pod> -- nc -z $DB_HOST $DB_PORT`
3. Check Redis: `redis-cli -h $REDIS_HOST ping`
4. If DB unreachable → failover to read replica; page DBA team.
5. If Redis unreachable → application degrades gracefully (local JWT cache, no dist lock); monitor for duplicate cron runs.

### `HighP99Latency` (> 2 s)

1. Pull traces from Jaeger for slow requests (filter `http.status_code=200`, sort by duration).
2. Check for N+1 patterns: `EXPLAIN ANALYZE` on queries from slow trace spans.
3. Check Redis latency: `redis-cli --latency-history -h $REDIS_HOST`
4. Scale horizontally if CPU-bound: `kubectl scale deployment gin-vue-admin --replicas=N`

### JWT Blacklist Miss (users not logged out)

1. Verify Redis connection: `redis-cli keys "gva:bl:*" | wc -l`
2. If Redis down, local in-memory blacklist is active — affected pods only; restart triggers clean state.
3. Emergency: rotate `GVA_JWT_SIGNING_KEY` to invalidate **all** existing tokens immediately.

### Casbin Policy Audit Chain Broken

1. Run chain verification:
   ```sql
   SELECT id, hash, prev_hash FROM sys_policy_change_logs ORDER BY id;
   ```
   Re-compute SHA-256 for each row using the same fields as `hashChainEntry()` in `policy_audit.go`.
2. A mismatch indicates tampering or a bug. Preserve all DB state; escalate to security team immediately.
3. Do **not** delete or update any rows — the table is append-only by design.

---

## Deployment & Rollback

### Standard Deploy

```bash
# 1. Build
docker build -t gin-vue-admin:$VERSION .

# 2. Tag & push
docker push registry.example.com/gin-vue-admin:$VERSION

# 3. Deploy (zero-downtime rolling)
kubectl set image deployment/gin-vue-admin app=registry.example.com/gin-vue-admin:$VERSION

# 4. Watch rollout
kubectl rollout status deployment/gin-vue-admin
```

### Rollback

```bash
# Immediate rollback to previous revision
kubectl rollout undo deployment/gin-vue-admin

# Rollback to specific revision
kubectl rollout undo deployment/gin-vue-admin --to-revision=N
```

**DB migrations are forward-only.** If a migration must be reversed, write a new migration — never run `DROP COLUMN` manually in production without DBA sign-off and a maintenance window.

---

## Disaster Recovery

### RTO / RPO Targets

| Scenario | RTO | RPO |
|----------|-----|-----|
| Single pod crash | < 1 min (k8s restart) | 0 (stateless) |
| DB primary failure | < 5 min (failover) | ≤ WAL lag (< 30 s typical) |
| Redis failure | < 2 min (restart or failover) | Acceptable — JWT blacklist rebuilds from DB |
| Full region outage | < 30 min (DR region activate) | ≤ last DB replica lag |

### DR Activation Steps

1. Promote DB read replica to primary.
2. Update `GVA_DB_DSN` secret in DR namespace.
3. Scale up DR deployment to normal replica count.
4. Update DNS / load balancer to point to DR region.
5. Verify health: `curl https://dr.example.com/health`
6. Notify stakeholders and open incident channel.

---

## Database Procedures

### Schema Migration

Migrations run automatically on startup via GORM `AutoMigrate`. For breaking changes (column drops, renames):

1. Deploy application version that writes both old and new columns (expand).
2. Backfill data; verify no reads of old column.
3. Deploy version that reads only new column.
4. Remove old column in a follow-up migration (contract).

### Clearing Old Records (Automated)

The `ClearTable` cron job runs daily and deletes:
- `sys_operation_records` older than 90 days
- `jwt_blacklists` older than 7 days

Whitelist is enforced in code (`clearTableWhitelist`); no other tables can be targeted.

---

## Secret Rotation

### JWT Signing Key

Rotating the key **immediately invalidates all active sessions** — plan for user re-login.

1. Generate new key: `openssl rand -hex 32`
2. Update secret: `kubectl create secret generic gva-secrets --from-literal=GVA_JWT_SIGNING_KEY=<new> --dry-run=client -o yaml | kubectl apply -f -`
3. Rolling restart: `kubectl rollout restart deployment/gin-vue-admin`
4. Verify no startup errors in logs.

### DB Password

1. Update password in DB.
2. Update `GVA_DB_PASSWORD` in secret store.
3. Rolling restart pods.

---

## Audit Log Verification

The `sys_policy_change_logs` table maintains a SHA-256 hash chain. To verify integrity:

```go
// Pseudocode — run as a scheduled job or pre-release check
rows = SELECT * FROM sys_policy_change_logs ORDER BY id ASC
prev = ""
for row in rows:
    expected = sha256(row.actor + "|" + row.action + ... + "|" + prev)
    if expected != row.hash:
        ALERT("chain broken at id=" + row.id)
    prev = row.hash
```

A broken chain means either:
- A row was modified after insert (tampering), or
- A bug in the application wrote an incorrect hash.

In either case: preserve DB state, escalate to security, do not modify rows.
