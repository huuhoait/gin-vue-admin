# Service Level Objectives (SLO)

> Revision: 2026-04-17 | Owner: Platform Engineering

## Definitions

| Term | Meaning |
|------|---------|
| **SLI** | Service Level Indicator — the measured metric |
| **SLO** | Service Level Objective — the target threshold |
| **Error Budget** | 100% − SLO; the allowable failure rate per rolling 30-day window |

---

## SLO Table

| # | SLI | SLO Target | Error Budget (30d) | Measurement Window |
|---|-----|-----------|--------------------|--------------------|
| 1 | API availability (`2xx + 3xx / total requests`) | ≥ 99.5% | 3h 36m downtime | Rolling 30 days |
| 2 | p95 response latency (all `/api/v1/*` endpoints) | ≤ 500 ms | — | Rolling 1 hour |
| 3 | p99 response latency | ≤ 2 000 ms | — | Rolling 1 hour |
| 4 | Login endpoint p95 latency | ≤ 800 ms | — | Rolling 1 hour |
| 5 | Authentication error rate (`401/403 / total auth requests`) | ≤ 0.1% | — | Rolling 24 hours |
| 6 | Background job success rate (cron tasks) | ≥ 99% | — | Per calendar day |
| 7 | Casbin policy propagation lag (change → enforcement) | ≤ 5 s | — | Per event |

---

## Error Budget Policy

- **> 50% budget consumed** in a rolling 30-day window → freeze all non-critical feature work; oncall must investigate.
- **> 75% budget consumed** → incident declared; postmortem required within 5 business days.
- **Budget exhausted** → SLO breach reported to stakeholders; release gate triggered.

---

## Alerting Thresholds (Prometheus / Grafana)

```yaml
# Availability — 1-hour burn rate alert (fires when on track to exhaust 30-day budget in < 6 h)
- alert: HighErrorRate
  expr: sum(rate(http_requests_total{status=~"5.."}[1h])) / sum(rate(http_requests_total[1h])) > 0.005
  for: 5m
  labels: { severity: warning }

# p95 latency breach
- alert: HighP95Latency
  expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 0.5
  for: 10m
  labels: { severity: warning }

# p99 latency breach
- alert: HighP99Latency
  expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 2.0
  for: 10m
  labels: { severity: critical }
```

---

## Instrumentation Notes

- Latency histogram: exported via OpenTelemetry SDK (see `core/otel.go`); collector scrapes OTLP gRPC on port 4317.
- Availability: computed from `http_requests_total` counter emitted by `otelgin` middleware.
- Traces: sampled at 10% in production (`GIN_MODE=release`), 100% in dev/staging.
