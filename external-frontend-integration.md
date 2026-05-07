# SkyAgent — External Frontend Integration Guide

> **Version:** 1.2.0 | **Last updated:** 2026-04-19
> **Applies to:** Vue 3 + Element Plus + axios clients (admin, agent portal, mobile web)
> **Backend reference:** Core (port 8001), Order (port 8002) — gọi **trực tiếp**, không qua API Gateway

Tài liệu duy nhất cho FE dev tích hợp **một dự án Vue admin tách rời** với SkyAgent backend. Bao gồm: kiến trúc, auth, envelope contract, pagination, error catalog, idempotency, PII rules, endpoint inventory, axios setup, i18n, **DTO reference đầy đủ** (§14), **full error code catalog** (§5).

---

## 1. Kiến trúc tích hợp

Theo mô hình **Epic 9 — BFF Proxy** (`_bmad-output/planning-artifacts/stories-epic-9.md` §9.1): FE chỉ giao tiếp với **Admin BFF (gin-vue-admin)**; BFF proxy forward request xuống Core/Order qua HTTP client tại `services/admin/server/plugin/skyagent/service/`. **API Gateway (:8000) KHÔNG được sử dụng** — BFF gọi trực tiếp Core và Order.

```
┌─────────────┐         ┌────────────────────┐         ┌─────────────┐
│  Vue Admin  │─HTTP──▶│ Admin BFF          │─HTTP──▶│  Core :8001 │
│  (Element   │         │ (gin-vue-admin)    │         ├─────────────┤
│   Plus)     │         │ Casbin RBAC +      │─HTTP──▶│  Order :8002│
└─────────────┘         │ proxy layer        │         └─────────────┘
                        │ maker/checker inj. │
                        └────────────────────┘
```

**Trách nhiệm mỗi tầng:**

| Tầng | Trách nhiệm |
|---|---|
| Vue Admin | UI, form validation phía client, gọi BFF qua axios |
| Admin BFF (gin-vue-admin) | JWT verify (shared `skyagent/shared/auth`), Casbin RBAC, inject `MakerID`/`CheckerID` từ JWT, proxy Core/Order, parse GVA envelope, timeout 10s, **không retry mutation** |
| Core `:8001` | Module `agent/`, `onboarding/`, `catalog/` + gRPC `:9901` |
| Order `:8002` | Module `order/`, `payment/`, `commission/`, webhooks |

**FE routing rule (qua BFF):**
- BFF expose lại các path `/v1/...` giống Core/Order để FE không phải biết service nào đứng sau.
- BFF nội bộ route theo prefix:
  - `/v1/agents/*`, `/v1/onboarding/*`, `/v1/catalog/*`, `/v1/documents/*` → Core
  - `/v1/orders/*`, `/v1/commission/*` → Order
  - `/v1/webhooks/*` → Order (BE-to-BE, FE không gọi)

**Tất cả HTTP endpoints có prefix `/v1`** (xem `services/core/cmd/server/main.go:191`).

---

## 2. Authentication

### Flow

```
1. POST /v1/auth/login         → trả { access_token, refresh_token, expires_in }
2. (per request) Authorization: Bearer <access_token>
3. Khi 401 + responseCode 1003 (TOKEN_EXPIRED) → POST /v1/auth/refresh
4. Logout → POST /v1/auth/logout (revoke refresh token)
```

### JWT claims (xem `shared/pkg/auth/claims.go`)

```ts
interface Claims {
  agent_id: string;     // UUID — chủ thể request
  agent_path: string;   // ltree path, ví dụ "1.5.23" — dùng cho ancestor check
  role: string;         // "agent" | "master_agent" | "admin" | ...
  permissions: string[];// ["agent.read", "agent.write", "order.refund", ...]
  exp: number;          // Unix epoch giây
  iat: number;
}
```

FE có thể decode JWT (không verify) để hide/show menu trước khi BE enforce. **Source of truth vẫn là BE** — BE re-verify mọi request qua middleware `JWTMiddleware(jwksURL)` (xem `shared/pkg/auth/middleware.go`).

### Headers bắt buộc

| Header | Khi nào | Ví dụ |
|---|---|---|
| `Authorization: Bearer <jwt>` | Mọi request trừ login/health/JWKS | `Bearer eyJhbGc...` |
| `X-Idempotency-Key: <uuid v4>` | Mọi `POST/PUT/PATCH/DELETE` mutation | `b3f9...` (xem §6) |
| `X-Trace-Id: <uuid>` | Optional — FE gen để correlate log | `a1b2...` |
| `Accept-Language: vi,en;q=0.9` | Optional — i18n hint | `vi` |

### JWKS verify (info để debug, FE không tự chạy)

```
GET /.well-known/jwks.json   (Admin BFF host)
→ { keys: [{ kty:"RSA", kid:"...", n:"...", e:"AQAB", use:"sig", alg:"RS256" }] }
```

BE verify bằng `kid` trong JWT header → match với `keys[].kid`.

---

## 3. Response Envelope

**Ground truth:** `shared/pkg/responsegva/response.go`

### Success

```json
{
  "code": 0,
  "data": { /* payload tuỳ endpoint */ },
  "msg": "Success"
}
```

`code: 0` ≡ thành công. `data` có thể là object, array, hoặc `null`.

### Error (chuẩn)

```json
{
  "code": 1001,
  "data": {
    "responseCode": 1001,
    "traceId": "a1b2c3...",
    "details": []
  },
  "msg": "Request lacks valid authentication credentials."
}
```

Khi `code !== 0`:
- HTTP status reflect mức nghiêm trọng (`401`, `403`, `404`, `422`, `429`, `500`...)
- `data` là `ErrorPayload` (xem `shared/pkg/responsegva/error.go`):
  ```ts
  interface ErrorPayload {
    responseCode: number;     // === envelope.code; map sang i18n key
    traceId?: string;         // dùng để báo support
    details?: FieldError[];   // có khi 422 validation
  }
  interface FieldError {
    field: string;            // tên field theo JSON, vd "phone"
    message: string;          // tiếng Anh, FE map sang vi qua i18n
  }
  ```

### Error 422 với field-level details (validation)

```json
{
  "code": 14001,
  "data": {
    "responseCode": 14001,
    "traceId": "a1b2...",
    "details": [
      { "field": "phone", "message": "must be a valid Vietnamese mobile number (10 digits, starting with 03/05/07/08/09)" },
      { "field": "email", "message": "invalid email format" }
    ]
  },
  "msg": "one or more fields failed validation"
}
```

FE bind `details[]` về form errors theo `field` name.

### HTTP status × envelope code matrix

| HTTP | Envelope `code` examples | Semantic | Recommended UI action |
|---|---|---|---|
| 200 | `0` | Success | render data |
| 201 | `0` | Created | toast + refresh |
| 202 | `0` | Accepted (async/Temporal) | poll workflow status |
| 400 | `3002`, `4001`, `4002`, `4003` | Bad request / missing idempotency | fix input, không retry |
| 401 | `1001`, `1003`, `2001-2005` | Unauthorized / token expired / signature err | refresh token (1003) hoặc redirect login |
| 403 | `1002` | Forbidden | hide action, toast "không đủ quyền" |
| 404 | `11001`, `11002` | Not found | toast / 404 page |
| 422 | `1xxxx-9xxxx` (domain validation) | Business rule violation | hiển thị `details[]` per field hoặc form-level toast |
| 429 | `3001` | Rate limited | exponential backoff |
| 500 | `5001`, `11007` | Internal error | toast + show traceId, encourage retry |
| 502/503/504 | `5002`, `5003`, `5004` | Upstream / timeout | auto retry với backoff |

---

## 4. Pagination

**Backend hỗ trợ 2 patterns. Endpoint quyết định pattern.** Đọc spec từng endpoint hoặc swagger.

### 4.1 Page-based (legacy / admin views)

Ground truth: `shared/pkg/responsegva/page_result.go`

**Request:** `?page=<n>&pageSize=<n>` (clamp: `page≥1`, `1≤pageSize≤100`)
**Response data:**
```json
{
  "list": [ /* items */ ],
  "total": 1234,
  "page": 1,
  "pageSize": 20
}
```

Use case: admin tables cần `total` cho pager component (`el-pagination`).

### 4.2 Cursor-based (preferred per CLAUDE.md R09)

**Request:** `?cursor=<opaque>&limit=<n>` (lần đầu omit `cursor`)
**Response data:**
```json
{
  "list": [ /* items */ ],
  "cursor": "eyJjcmVhdGVkX2F0...",
  "has_more": true
}
```

Use case: feeds, infinite scroll, large dataset. Không có `total` (perf).

**FE nhận `cursor=null` hoặc `has_more=false` ⇒ stop pagination.**

---

## 5. Error Codes Catalog

**Code format:** `<domain><category><sequence>` — vd `11005` = domain 1 (Agent) / category 1 (Hierarchy) / sequence 005.

**Source:** `shared/pkg/rfc7807/codes/*.go`. **Tổng hiện tại: 78 codes** (xác nhận bằng `grep -rh "rfc7807.Register" shared/pkg/rfc7807/codes/`).

### Domain map

| Domain | Range | File | Số lượng |
|---|---|---|---|
| 0 — Common (auth, sig, rate limit, validation, infra) | 1xxx–5xxx | `common.go` | 17 |
| 1 — Agent | 1xxxx | `agent.go` | 17 |
| 2 — Onboarding/KYC | 2xxxx | `onboarding.go` | 12 |
| 3 — Order | 3xxxx | `order.go` | 12 |
| 4 — Commission | 4xxxx | `commission.go` | 4 |
| 5 — Supplier | 5xxxx | `supplier.go` | 3 |
| 6 — Catalog | 6xxxx | `catalog.go` | 8 |
| 7 — Report | 7xxxx | `report.go` | 3 |

### 5.1 Common (auth, signature, rate limit, validation, infra)

| Code | Constant | HTTP | Retryable | Title / UI action |
|---|---|---|---|---|
| `1001` | `ErrUnauthorized` | 401 | no | Unauthorized → redirect login |
| `1002` | `ErrForbidden` | 403 | no | Forbidden → toast "không đủ quyền" |
| `1003` | `ErrTokenExpired` | 401 | **yes** | Token expired → refresh + replay |
| `2001` | `ErrSignatureMissing` | 401 | no | RFC 9421 sig required |
| `2002` | `ErrSignatureInvalid` | 401 | no | Sig verify failed |
| `2003` | `ErrSignatureExpired` | 401 | **yes** | Sig timestamp > 5min — resign |
| `2004` | `ErrSignatureReplay` | 401 | no | Nonce reused — gen new |
| `2005` | `ErrSignatureKeyRevoked` | 401 | no | Reissue key pair |
| `3001` | `ErrRateLimited` | 429 | **yes** | Backoff + retry |
| `3002` | `ErrIdempotencyKeyRequired` | 400 | no | Gen UUID → retry |
| `4001` | `ErrInvalidRequestBody` | 400 | no | Check JSON syntax |
| `4002` | `ErrInvalidQueryParam` | 400 | no | Check query format |
| `4003` | `ErrInvalidPathParam` | 400 | no | Check path UUID |
| `5001` | `ErrInternalError` | 500 | **yes** | Toast + traceId |
| `5002` | `ErrServiceUnavailable` | 503 | **yes** | Retry với backoff |
| `5003` | `ErrRequestTimeout` | 504 | **yes** | Retry |
| `5004` | `ErrDependencyFailure` | 502 | **yes** | Upstream lỗi — retry |

### 5.2 Agent (domain 1)

| Code | Constant | HTTP | Retryable | UI action |
|---|---|---|---|---|
| `11001` | `ErrAgentNotFound` | 404 | no | Check ID / 404 page |
| `11002` | `ErrAgentParentNotFound` | 404 | no | Check parent_id |
| `11003` | `ErrAgentParentNotActive` | 422 | no | Parent phải ACTIVE |
| `11004` | `ErrAgentParentMaxChildren` | 422 | no | Parent đã đạt max con |
| `11005` | `ErrAgentMaxDepthExceeded` | 422 | no | Vượt depth cho phép |
| `11006` | `ErrAgentInvalidLevel` | 422 | no | Level không hợp lệ |
| `11007` | `ErrAgentSubtreeCascadeFailed` | 500 | **yes** | Cascade fail — retry hoặc support |
| `12001` | `ErrAgentNotActive` | 422 | no | Kích hoạt agent trước |
| `12002` | `ErrAgentPendingApproval` | 422 | wait | Chờ duyệt |
| `12003` | `ErrAgentInvalidStatusTransition` | 422 | no | Check status hiện tại |
| `12004` | `ErrAgentTerminated` | 422 | no | Agent đã chấm dứt |
| `13001` | `ErrAgentDuplicatePhone` | 409 | no | Phone đã tồn tại |
| `13002` | `ErrAgentDuplicateEmail` | 409 | no | Email đã tồn tại |
| `13003` | `ErrAgentInvalidPhone` | 400 | no | Sai format VN phone |
| `13004` | `ErrAgentInvalidEmail` | 400 | no | Sai format email |
| `14001` | `ErrAgentValidationFailed` | 400 | no | Show `details[]` |
| `15001` | `ErrAgentConcurrentModification` | 409 | **yes** | Re-fetch + retry |

### 5.3 Onboarding / KYC (domain 2)

| Code | Constant | HTTP | Retryable | UI action |
|---|---|---|---|---|
| `21001` | `ErrKYCAlreadyPending` | 409 | no | KYC đã pending |
| `21002` | `ErrKYCDuplicateIDNumber` | 409 | no | CCCD đã dùng |
| `21003` | `ErrKYCProviderUnavailable` | 502 | **yes** | Provider down — retry |
| `21004` | `ErrKYCDocumentExpired` | 422 | no | Giấy tờ hết hạn |
| `21005` | `ErrKYCAgentNotEligible` | 422 | no | Agent không đủ điều kiện |
| `21006` | `ErrKYCNotInManualReview` | 422 | no | Không ở trạng thái review |
| `23001` | `ErrKYCSelfReview` | 403 | no | Không được tự duyệt |
| `23003` | `ErrKYCRejectNoteRequired` | 400 | no | Nhập lý do reject |
| `24001` | `ErrKYCSubmissionNotFound` | 404 | no | Check submission_id |
| `25001` | `ErrKYCInvalidFileType` | 400 | no | Chỉ jpg/png/heic |
| `25002` | `ErrKYCFileTooLarge` | 400 | no | Tối đa 10MB |
| `25003` | `ErrKYCInvalidIDNumber` | 400 | no | CCCD 12 chữ số |

### 5.4 Order (domain 3)

| Code | Constant | HTTP | Retryable | UI action |
|---|---|---|---|---|
| `31001` | `ErrOrderNotFound` | 404 | no | Check order_id |
| `31002` | `ErrOrderAgentNotActive` | 422 | no | Agent không ACTIVE |
| `31003` | `ErrOrderAlreadyCancelled` | 409 | no | Đã cancel |
| `31004` | `ErrOrderInvalidPaymentMethod` | 400 | no | Chọn UPC / Cash App / vnpay / momo |
| `31005` | `ErrOrderRefundNotAllowed` | 422 | no | Không refund ở trạng thái này |
| `31006` | `ErrOrderRefundAmountExceeded` | 400 | no | Amount > order total |
| `31007` | `ErrOrderAccessDenied` | 403 | no | Không có quyền xem order |
| `32001` | `ErrOrderPaymentFailed` | 422 | **yes** | Payment fail — retry |
| `32002` | `ErrVNPayInvalidSignature` | 403 | no | VNPay sig sai (webhook) |
| `32003` | `ErrVNPayAmountMismatch` | 400 | no | VNPay amount lệch |
| `32004` | `ErrMoMoInvalidSignature` | 403 | no | MoMo sig sai |
| `32005` | `ErrMoMoAmountMismatch` | 400 | no | MoMo amount lệch |

### 5.5 Commission (domain 4)

| Code | Constant | HTTP | Retryable | UI action |
|---|---|---|---|---|
| `41001` | `ErrCascadeDuplicate` | 409 | no | Cascade rule đã tồn tại |
| `41002` | `ErrCommissionConfigNotFound` | 404 | no | Không có cấu hình |
| `42001` | `ErrLedgerImbalanced` | 500 | no | Ledger lệch — gọi support |
| `42002` | `ErrLedgerEntryNotFound` | 404 | no | Entry không tồn tại |

### 5.6 Supplier (domain 5)

| Code | Constant | HTTP | Retryable | UI action |
|---|---|---|---|---|
| `51001` | `ErrSupplierNotFound` | 404 | no | Check supplier code |
| `52001` | `ErrSupplierAPIError` | 502 | **yes** | Supplier lỗi — retry |
| `52002` | `ErrSupplierCircuitOpen` | 503 | **yes** | Circuit breaker — chờ |

### 5.7 Catalog (domain 6)

| Code | Constant | HTTP | Retryable | UI action |
|---|---|---|---|---|
| `61001` | `ErrCategoryNotFound` | 404 | no | Check category_id |
| `62001` | `ErrCatalogAccessDenied` | 403 | no | Không có quyền |
| `62002` | `ErrKYCTierInsufficient` | 403 | no | Nâng KYC tier |
| `62003` | `ErrDisabledByParent` | 403 | no | Parent đã disable |
| `62004` | `ErrAncestorSuspended` | 403 | no | Tổ tiên bị suspend |
| `62005` | `ErrInsuranceCertRequired` | 403 | no | Cần chứng chỉ bảo hiểm (TT 50/2023) |
| `63001` | `ErrCommissionRateInvalid` | 400 | no | Rate phải ≥ 0 |
| `63002` | `ErrCommissionHardCapExceeded` | 422 | no | Vượt cap cấu hình |

### 5.8 Report (domain 7)

| Code | Constant | HTTP | Retryable | UI action |
|---|---|---|---|---|
| `71001` | `ErrReportNotFound` | 404 | no | Check report_id |
| `72001` | `ErrReconAlreadyRunning` | 409 | no | Recon đang chạy — chờ |
| `72002` | `ErrExportInProgress` | 409 | no | Export đang chạy — chờ |

### Retryable flag cheat-sheet

FE nên tự retry (với exponential backoff) cho các code: `1003`, `2003`, `3001`, `5001`, `5002`, `5003`, `5004`, `11007`, `15001`, `21003`, `32001`, `52001`, `52002`. Các code 422 không bao giờ retry (client phải sửa input). Các code 409 thường không retry, trừ `15001` (concurrent modification → re-fetch rồi retry).

### i18n key convention

```
errors.<responseCode>           // Ví dụ: errors.11001 = "Không tìm thấy đại lý"
errors.unknown                  // Fallback khi code chưa map
```

Khi BE thêm code mới, FE phải thêm 2 key (vi + en). Xem §10.

---

## 6. Idempotency

**Per CLAUDE.md R06 + code `3002`:** mọi mutation (`POST/PUT/PATCH/DELETE` thay đổi state) **bắt buộc** header:

```
X-Idempotency-Key: <uuid v4>
```

BE dedupe bằng `event_id` + Redis TTL (24h). Cùng key + cùng body → trả về kết quả cũ. Cùng key + body khác → 409 conflict.

### Generate ở FE

```ts
import { v4 as uuid } from 'uuid';

axios.interceptors.request.use((cfg) => {
  if (['post', 'put', 'patch', 'delete'].includes(cfg.method?.toLowerCase() ?? '')) {
    cfg.headers ??= {};
    if (!cfg.headers['X-Idempotency-Key']) {
      cfg.headers['X-Idempotency-Key'] = uuid();
    }
  }
  return cfg;
});
```

**Khi nào FE pass `X-Idempotency-Key` thủ công:**
- Submit form mà user có thể click double — gen 1 key trước, reuse cho mọi retry trong 24h
- Webhook resend / replay scenarios

---

## 7. PII & Response Masking

**Per CLAUDE.md R04 + Core CLAUDE.md §5:** các field PII được mã hoá AES-256-GCM trong DB; response trả về **đã mask**:

| Field | Format mask | Ví dụ |
|---|---|---|
| `phone` / `contact_phone` | `<3 đầu>****<3 cuối>` | `090****890` |
| `email` / `contact_email` | `<3 đầu>****@<domain>` | `tha****@hdbank.com.vn` |
| `cccd` | `<3 đầu>******<3 cuối>` | `036******123` |
| `bank_account` | `****<4 cuối>` | `****1234` |
| `representative_cccd` | same as `cccd` | — |
| `id` (đôi khi) | `<5 đầu>...<3 cuối>` | `b3f9a...123` |

**FE LƯU Ý:**
- KHÔNG attempt unmask hoặc gọi API "get full PII" trừ khi role=admin và endpoint admin riêng (`/v1/admin/agents/:id`)
- Search dùng plaintext tìm sẽ thất bại — BE dùng blind index HMAC; FE chỉ truyền plaintext, BE hash trước khi query
- Hiển thị mask trên UI là **mặc định**; expose full chỉ qua "Show Full" button (nếu role allow)

---

## 8. Endpoint Inventory

> **Lưu ý:** Bảng dưới là baseline; nguồn chính xác là Swagger (xem §9). Path `:id` = UUID.

### 8.1 Agent (`services/core/internal/agent/handler/http/`)

| Method | Path | Mô tả | Notes |
|---|---|---|---|
| `POST` | `/v1/agents` | Tạo đại lý mới (pending) | Idempotency required |
| `POST` | `/v1/agents/:id/children` | Tạo sub-agent dưới parent | Async qua Temporal nếu enabled, fallback sync |
| `GET` | `/v1/agents` | List với filter (`status`, `level`, `parent_id`) | Cursor pagination |
| `GET` | `/v1/agents/:id` | Detail (PII masked) | — |
| `GET` | `/v1/agents/:id/tree` | Subtree (depth limit) | Cursor pagination |
| `GET` | `/v1/agents/:id/ancestors` | Chain tổ tiên (root → :id) | — |
| `PUT` | `/v1/agents/:id` | Update non-PII fields | — |
| `PUT` | `/v1/agents/:id/status` | Suspend / Reactivate / Terminate | Body: `{ action, reason }` |
| `POST` | `/v1/agents/:id/documents` | Upload document (KYC, bank statement, ...) | Multipart |
| `GET` | `/v1/agents/:id/documents` | List documents | — |
| `PUT` | `/v1/agents/:id/bank` | Update bank info | — |
| `PUT` | `/v1/documents/:id/review` | Reviewer approve/reject doc | Reviewer-only |
| `GET` | `/v1/admin/agents/:id` | Admin view (full PII, role=admin) | Admin-only |

### 8.2 Onboarding/KYC (`services/core/internal/onboarding/handler/http/`)

Endpoint chính (xem source để full list):
- `POST /v1/onboarding/kyc/submit` — multipart (CCCD front/back, selfie); idempotency required
- `GET /v1/onboarding/kyc/pending` — reviewer queue
- `GET /v1/onboarding/kyc/:id` — detail
- `POST /v1/onboarding/kyc/:id/approve` — reviewer approve
- `POST /v1/onboarding/kyc/:id/reject` — reviewer reject
- `POST /v1/onboarding/agents` — onboarding ticket (Epic 12 — full agent_payload JSON)
- `POST /v1/onboarding/tickets/:id/promote` — promote ticket sang `agents` table

### 8.3 Catalog (`services/core/internal/catalog/handler/http/`)

- `GET /v1/catalog/categories` — list categories với access status
- `GET /v1/catalog/categories/:id/products` — products theo category
- `POST /v1/catalog/products/:id/visibility` — parent toggle on/off cho con
- `PUT /v1/catalog/agents/:id/commission-overrides` — commission per agent × category

`CheckProductAccess` là **gRPC** (`proto/catalog/v1/`), không HTTP. FE thường không gọi.

### 8.4 Order (`services/order/internal/order/handler/http/`)

- `POST /v1/orders` — tạo đơn (idempotency required)
- `GET /v1/orders` — list (cursor)
- `GET /v1/orders/:id` — detail
- `POST /v1/orders/:id/cancel` — cancel
- `POST /v1/orders/:id/refund` — refund full/partial
- `GET /v1/orders/:id/status-logs` — audit trail

**Webhooks (BE-to-BE, FE chỉ document để biết shape):**
- `POST /v1/webhooks/payment/:gateway` (`vnpay`, `momo`)
- `POST /v1/webhooks/supplier/:code`

### 8.5 Commission

- `GET /v1/commission/wallet` — `{ pending, available, frozen, debt }` của user hiện tại
- `GET /v1/commission/preview?order_id=...` — cascade 5 levels
- `GET /v1/commission/entries` — list (filter: status, order_id, ...)
- `POST /v1/commission/withdraw` — yêu cầu rút (FR-C13)
- `GET /v1/commission/withdrawals` — history

---

## 9. Swagger / OpenAPI for Codegen

### Generate per service

```bash
cd services/core  && make swagger    # → services/core/docs/swagger.{json,yaml}
cd services/order && make swagger    # → services/order/docs/swagger.{json,yaml}
```

### Aggregate cho FE

Khuyến nghị tạo script `scripts/aggregate-openapi.sh` (xem _bmad-output/implementation-artifacts/15-5):

```bash
#!/usr/bin/env bash
set -euo pipefail
mkdir -p docs/openapi
for svc in core order; do
  (cd services/$svc && make swagger)
  cp services/$svc/docs/swagger.json docs/openapi/$svc.json
done
```

### TypeScript codegen (FE side)

```bash
# Option 1: openapi-typescript (chỉ types, nhẹ)
npx openapi-typescript ./docs/openapi/core.json -o src/api/types/core.ts

# Option 2: openapi-generator-cli (full SDK với axios client)
npx @openapitools/openapi-generator-cli generate \
  -i ./docs/openapi/core.json \
  -g typescript-axios \
  -o src/api/sdk/core
```

### Apifox import

UI: Project → Import → OpenAPI → upload `docs/openapi/core.json` (hoặc URL nếu publish).

---

## 10. axios Setup (TypeScript, copy-paste)

### `src/api/contract.ts` — types + ApiError

```ts
export interface SuccessEnvelope<T> { code: 0; data: T; msg: string }

export interface FieldError { field: string; message: string }
export interface ErrorPayload {
  responseCode: number;
  traceId?: string;
  details?: FieldError[];
}
export interface ErrorEnvelope { code: number; data: ErrorPayload; msg: string }

export interface PageResult<T> { list: T[]; total: number; page: number; pageSize: number }
export interface CursorPage<T> { list: T[]; cursor: string | null; has_more: boolean }

export class ApiError extends Error {
  constructor(
    public responseCode: number,
    public traceId: string | undefined,
    public details: FieldError[],
    public httpStatus: number,
    msg: string,
  ) { super(msg); this.name = 'ApiError'; }
}
```

### `src/api/client.ts` — single axios instance

```ts
import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { v4 as uuid } from 'uuid';
import { ApiError, type ErrorEnvelope, type SuccessEnvelope } from './contract';
import { useAuthStore } from '@/stores/auth';
import { i18n } from '@/i18n';

// FE chỉ gọi Admin BFF (gin-vue-admin); BFF proxy xuống Core/Order.
// Gateway :8000 KHÔNG dùng — xem §1.
export const api = axios.create({
  baseURL: import.meta.env.VITE_ADMIN_BFF_BASE ?? 'http://localhost:8888',
  timeout: 15_000,
});

// Request interceptor — auth + idempotency
api.interceptors.request.use((cfg: InternalAxiosRequestConfig) => {
  const token = useAuthStore().accessToken;
  if (token) cfg.headers.Authorization = `Bearer ${token}`;

  const method = (cfg.method ?? 'get').toLowerCase();
  if (['post', 'put', 'patch', 'delete'].includes(method)) {
    cfg.headers['X-Idempotency-Key'] ??= uuid();
  }

  cfg.headers['Accept-Language'] = i18n.global.locale.value;
  return cfg;
});

// Response interceptor — unwrap envelope + handle 401 refresh
let isRefreshing = false;
let pendingQueue: Array<{ resolve: (v: unknown) => void; reject: (e: unknown) => void; cfg: InternalAxiosRequestConfig }> = [];

instances.forEach((inst) => {
  inst.interceptors.response.use(
    (res) => {
      const env = res.data as SuccessEnvelope<unknown> | ErrorEnvelope;
      if (env.code === 0) return env.data; // ⬅ FE chỉ nhận `data`, không thấy envelope
      // Defensive: BE trả 200 nhưng code !== 0 (không nên xảy ra)
      throw mapError(res.status, env as ErrorEnvelope);
    },
    async (err: AxiosError) => {
      const status = err.response?.status ?? 0;
      const env = err.response?.data as ErrorEnvelope | undefined;

      // 401 + TOKEN_EXPIRED → try refresh, replay trên CÙNG instance
      if (status === 401 && env?.data?.responseCode === 1003) {
        const auth = useAuthStore();

        if (isRefreshing) {
          return new Promise((resolve, reject) => {
            pendingQueue.push({ resolve, reject, cfg: err.config! });
          });
        }

        isRefreshing = true;
        try {
          await auth.refresh();
          pendingQueue.forEach(({ resolve, cfg }) => {
            cfg.headers!.Authorization = `Bearer ${auth.accessToken}`;
            resolve(inst(cfg));
          });
          pendingQueue = [];
          err.config!.headers!.Authorization = `Bearer ${auth.accessToken}`;
          return inst(err.config!);
        } catch (refreshErr) {
          pendingQueue.forEach(({ reject }) => reject(refreshErr));
          pendingQueue = [];
          auth.logout();
          throw mapError(status, env);
        } finally {
          isRefreshing = false;
        }
      }

      throw mapError(status, env);
    },
  );
});

function mapError(httpStatus: number, env: ErrorEnvelope | undefined): ApiError {
  if (!env || typeof env.code !== 'number') {
    return new ApiError(5001, undefined, [], httpStatus, 'Network or unknown error');
  }
  return new ApiError(
    env.data?.responseCode ?? env.code,
    env.data?.traceId,
    env.data?.details ?? [],
    httpStatus,
    env.msg,
  );
}

export function extractFieldErrors(err: unknown): Record<string, string> {
  if (!(err instanceof ApiError)) return {};
  return Object.fromEntries(err.details.map((d) => [d.field, d.message]));
}

export function translateError(err: unknown): string {
  if (!(err instanceof ApiError)) return i18n.global.t('errors.unknown');
  const key = `errors.${err.responseCode}`;
  const t = i18n.global.t(key);
  return t === key ? (err.message || i18n.global.t('errors.unknown')) : t;
}
```

### `src/api/agent.ts` — example module

```ts
import { api } from './client';
import type { CursorPage } from './contract';

export interface Agent {
  id: string;
  agent_code: string;
  contact_phone: string;   // mask: 090****890
  contact_email: string;
  status: 'pending' | 'active' | 'suspended' | 'terminated';
  level: number;
  created_at: string;      // RFC3339
}

export async function listAgents(params: { cursor?: string; limit?: number; status?: string }) {
  return api.get<unknown, CursorPage<Agent>>('/v1/agents', { params });
}

export async function createAgent(payload: {
  parent_id: string;
  full_name: string;
  phone: string;
  email: string;
  province: string;
  district: string;
  business_type: string;
}) {
  return api.post<unknown, Agent>('/v1/agents', payload);
  // X-Idempotency-Key tự động gắn bởi interceptor
}
```

### Vue component example

```vue
<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { listAgents, createAgent, type Agent } from '@/api/agent';
import { extractFieldErrors, translateError, ApiError } from '@/api/client';
import { ref, reactive } from 'vue';

const agents = ref<Agent[]>([]);
const cursor = ref<string | null>(null);
const hasMore = ref(true);
const formErrors = reactive<Record<string, string>>({});

async function loadMore() {
  if (!hasMore.value) return;
  const page = await listAgents({ cursor: cursor.value ?? undefined, limit: 20 });
  agents.value.push(...page.list);
  cursor.value = page.cursor;
  hasMore.value = page.has_more;
}

async function onSubmit(payload: any) {
  Object.keys(formErrors).forEach((k) => delete formErrors[k]);
  try {
    const agent = await createAgent(payload);
    ElMessage.success(`Created ${agent.agent_code}`);
  } catch (e) {
    if (e instanceof ApiError && e.httpStatus === 422) {
      Object.assign(formErrors, extractFieldErrors(e));
    } else {
      ElMessage.error(translateError(e));
    }
  }
}
</script>
```

---

## 11. i18n Mapping

### Locale file structure (Vue admin project)

```jsonc
// src/locales/vi.json
{
  "errors": {
    "unknown": "Đã có lỗi xảy ra",
    "1001": "Yêu cầu chưa xác thực",
    "1002": "Bạn không có quyền thực hiện",
    "1003": "Phiên đăng nhập hết hạn",
    "3001": "Quá nhiều yêu cầu, thử lại sau",
    "3002": "Thiếu khóa idempotency",
    "11001": "Không tìm thấy đại lý",
    "11003": "Đại lý cha chưa được kích hoạt",
    "12001": "Đại lý chưa được kích hoạt",
    "14001": "Dữ liệu không hợp lệ"
  },
  "agent": {
    "create": "Tạo đại lý",
    "list": "Danh sách đại lý",
    "status": {
      "pending": "Chờ duyệt",
      "active": "Hoạt động",
      "suspended": "Tạm ngưng",
      "terminated": "Đã chấm dứt"
    }
  }
}
```

```jsonc
// src/locales/en.json
{
  "errors": {
    "unknown": "Unknown error",
    "1001": "Unauthorized",
    "1002": "Forbidden",
    "1003": "Session expired",
    "3001": "Too many requests",
    "3002": "Idempotency key required",
    "11001": "Agent not found",
    "11003": "Parent agent not active",
    "12001": "Agent not active",
    "14001": "Validation failed"
  }
}
```

### Convention thêm key mới

Khi BE merge thêm code (ví dụ `5005 NEW_DEPENDENCY_FAILURE`):

1. Thêm `errors.5005` vào `vi.json`
2. Thêm `errors.5005` vào `en.json`
3. (Optional) Note vào `docs/openapi/` regen

Recommend setup script CI verify cả 2 locale có cùng set keys:

```bash
diff <(jq -r 'keys' vi.json | sort) <(jq -r 'keys' en.json | sort)
```

---

## 12. Testing & Debugging

### Smoke test với curl

> FE gọi qua Admin BFF (`:8888` — gin-vue-admin default). BE debug có thể gọi thẳng Core (`:8001`).

```bash
BFF=http://localhost:8888

# 1. Login qua BFF
TOKEN=$(curl -s -X POST $BFF/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' \
  | jq -r '.data.access_token')

# 2. List agents (BFF proxy → Core)
curl -s $BFF/v1/agents?limit=5 \
  -H "Authorization: Bearer $TOKEN" \
  | jq .

# Expected envelope:
# { "code": 0, "data": { "list": [...], "cursor": "...", "has_more": true }, "msg": "Success" }

# 3. Trigger validation error (POST với phone sai)
curl -s -X POST $BFF/v1/agents \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Idempotency-Key: $(uuidgen)" \
  -H "Content-Type: application/json" \
  -d '{"phone":"123","email":"bad","full_name":"Test","province":"HN","district":"BD","business_type":"individual"}' \
  | jq .

# Expected: code=14001, data.details[] có {field, message}

# 4. Bypass BFF để debug Core trực tiếp (BE-only, FE không làm)
# curl -s http://localhost:8001/v1/agents?limit=5 -H "Authorization: Bearer $TOKEN" | jq .
```

### Decode JWT (debug)

```bash
echo $TOKEN | cut -d'.' -f2 | base64 -d | jq .
# → { "agent_id": "...", "role": "admin", "permissions": [...], "exp": ..., "iat": ... }
```

### Common pitfalls

| Lỗi FE | Nguyên nhân | Cách sửa |
|---|---|---|
| `code: 3002` trên POST | Quên `X-Idempotency-Key` | Đảm bảo interceptor request gắn cho mọi mutation |
| `code: 1003` repeat sau khi refresh | Refresh không update store | Sau `auth.refresh()` phải `useAuthStore().setAccessToken(...)` |
| Field mask `090****890` không pass validation | FE gửi lại masked value khi update | Update form chỉ gửi fields user thật sự edit, không gửi lại masked PII |
| 422 nhưng `details` rỗng | Lỗi không phải binding validation | Check `responseCode` để biết business rule nào fail |
| 401 trên endpoint public | Client gắn header expired token | Login flow phải skip interceptor (whitelist `/v1/auth/login`, `/v1/health`) |

---

## 13. Versioning & Maintenance

### Spec drift detection

| Artifact | Owner | Update khi nào |
|---|---|---|
| `shared/pkg/responsegva/` | BE | Hiếm khi (envelope stable) |
| `shared/pkg/rfc7807/codes/*.go` | BE | Mỗi lần thêm error code mới |
| `services/{core,order}/docs/swagger.{json,yaml}` | BE | Auto-regen mỗi build (`make swagger`) |
| `docs/external-frontend-integration.md` | BE + FE shared | Khi: thêm endpoint mới, đổi auth flow, đổi envelope, thêm header convention |
| Vue admin project locale files | FE | Mỗi lần BE thêm error code |

### PR checklist (khuyến nghị thêm vào template)

- [ ] Nếu thêm endpoint mới → update §8 Endpoint Inventory + regen swagger
- [ ] Nếu thêm error code → update §5 Frequently-used codes (nếu phổ biến) + báo FE thêm i18n key
- [ ] Nếu đổi response shape → bump `Version` ở header doc + note breaking change
- [ ] Nếu đổi auth flow → update §2

### Reference

- `CLAUDE.md` (root) — toàn bộ rules R01-R16
- `services/core/CLAUDE.md` — Core-specific (PII, ltree, transaction)
- `shared/pkg/responsegva/` — envelope source
- `shared/pkg/rfc7807/codes/` — error catalog source
- `shared/pkg/auth/` — JWT/JWKS verify

---

## 14. DTO Reference (đầy đủ request/response shapes)

> **Ground truth:** `services/core/internal/*/handler/http/dto/` và `services/order/internal/*/handler/http/dto/`. Bảng dưới quote **y nguyên** JSON tag + binding tag từ struct Go.
>
> **Ký hiệu:** `?` = optional. **Custom validators:** `vn_phone` = `^(03|05|07|08|09)\d{8}$`. `status_change_reason` = required (1–500 ký tự) khi `action ∈ {suspend, terminate}`.

### 14.1 Agent module

#### `POST /v1/agents` — CreateAgentRequest → 201 CreateAgentResponse

**Request (JSON):**

| Field | Type | Binding | Notes |
|---|---|---|---|
| `parent_id` | string | `required,uuid` | UUID agent cha (L0 vẫn required nhưng có thể dùng "null-UUID" hoặc omit tùy BE convention) |
| `full_name` | string | `required,min=1,max=255` | |
| `phone` | string | `required,vn_phone` | 10 số VN |
| `email` | string | `required,email,max=320` | |
| `province` | string | `required,min=1,max=100` | |
| `district?` | string | `omitempty,max=100` | |
| `business_type?` | string | `omitempty,oneof=COMPANY HOUSEHOLD INDIVIDUAL` | |
| `address_line?` | string | `omitempty,max=500` | |
| `province_code?` | string | `omitempty,max=10` | |
| `district_code?` | string | `omitempty,max=10` | |
| `ward_code?` | string | `omitempty,max=10` | |
| `partner_id?` | string | `omitempty,max=100` | |
| `user_id?` | string | `omitempty,uuid` | |
| `referral_code?` | string | `omitempty,max=50` | |
| `referral_name?` | string | `omitempty,max=255` | |
| `maker_id?` | string | `omitempty,uuid` | BFF inject từ JWT |
| `checker_id?` | string | `omitempty,uuid` | BFF inject từ JWT |

**Response (201):**

```ts
interface CreateAgentResponse {
  id: string;           // agent UUID
  code: string;         // "L0-00001"
  level: number;        // 0–4
  path: string;         // ltree, vd "0.1.2"
  status: string;       // "pending_approval" cho L0, "active" cho sub-agent auto-approve
  kyc_tier: number;     // 0
  created_at: string;   // RFC3339
  created_by: string;
}
```

#### `POST /v1/agents/:id/children` — CreateSubAgentRequest → 201 CreateAgentResponse

**Request (JSON):**

```ts
interface CreateSubAgentRequest {
  identity: {
    full_name: string;    // required, 1–255
    phone: string;        // required, vn_phone
    email: string;        // required, email, max 320
    province: string;     // required, 1–100
    address?: string;     // omitempty, max 500 (chưa map vào agents schema)
    agent_name?: string;  // omitempty, max 255
    userId?: string;      // omitempty, uuid (camelCase — lưu ý!)
  };
}
```

#### `GET /v1/agents` — ListAgentsQueryParams → 200 `PageResult<AgentSummaryDTO>` hoặc `CursorPage<AgentSummaryDTO>`

**Query params:**

| Param | Type | Binding | Notes |
|---|---|---|---|
| `status?` | string | `oneof=active pending_approval suspended terminated` | |
| `province?` | string | `max=100` | |
| `level?` | int | `min=0,max=4` | |
| `parent_id?` | string | `uuid` | |
| `keyword?` | string | `max=255` | |
| `referral_code?` | string | `max=50` | |
| `kyc_tier?` | int | `min=0,max=2` | |
| `created_from?` | string | — | YYYY-MM-DD |
| `created_to?` | string | — | YYYY-MM-DD |
| `page?` | int | `min=1` | Offset pagination |
| `pageSize?` | int | `min=1,max=100` | Offset pagination |
| `cursor?` | string | `max=2048` | Cursor pagination (preferred) |

**Response item (`AgentSummaryDTO`):**

```ts
interface AgentSummaryDTO {
  id: string;
  code: string;
  full_name: string;       // plaintext trong list
  level: number;
  path: string;            // ltree
  status: string;
  kyc_tier: number;
  province: string;
  parent_id?: string;
  children_count: number;
  created_at: string;
}
```

#### `GET /v1/agents/:id` → 200 `AgentDetailDTO` (PII masked)

```ts
interface AgentDetailDTO extends AgentSummaryDTO {
  phone?: string;                // "****4567"
  email?: string;                // "a***@example.com"
  district?: string;
  avatar_url?: string;
  business_name?: string;
  business_type?: string;        // COMPANY|HOUSEHOLD|INDIVIDUAL
  tax_code?: string;
  representative_name?: string;
  representative_cccd?: string;  // "****8901"
  bank_name?: string;
  bank_account?: string;         // "****1234"
  permanent_address?: string;
  contact_address?: string;
  attachments: AttachmentDTO[];
  max_children: number;
  referral_code?: string;
  referral_name?: string;
  partner_id?: string;
  user_id?: string;
  updated_at: string;
  created_by: string;
  updated_by: string;
}

interface AttachmentDTO {
  type: string;        // BUSINESS_LICENSE | ID_CARD_FRONT | ...
  url: string;         // S3 pre-signed URL
  uploaded_at?: string;
}
```

#### `GET /v1/admin/agents/:id` → 200 `AgentAdminDetailDTO` (admin-only, plaintext PII)

Bằng `AgentDetailDTO` + thêm **4 trường plaintext** — **response có header `Cache-Control: no-store`**:

```ts
interface AgentAdminDetailDTO extends AgentDetailDTO {
  phone_plain?: string;
  email_plain?: string;
  representative_cccd_plain?: string;
  bank_account_plain?: string;
}
```

#### `GET /v1/agents/:id/tree` → 200 cursor list

**Query:** `max_depth?` (int, 1–5, default 1), `page_size?` (int, 1–200, default 50), `cursor?` (string, max 2048).
**Response:** `{ data: AgentSummaryDTO[], cursor: string|null, has_more: bool, total_count: number }`.

#### `PUT /v1/agents/:id/status` — UpdateStatusRequest

```ts
interface UpdateStatusRequest {
  action: 'approve' | 'suspend' | 'reactivate' | 'terminate';  // required
  reason?: string;        // required (1–500) khi action ∈ {suspend, terminate}
  maker_id?: string;      // uuid, BFF inject
  checker_id?: string;    // uuid, BFF inject
}
```

**Response shape tùy action:**
- `approve` → `{ id, status: "active", approved_by, approved_at }`
- `suspend` → `{ id, status: "suspended", affected_count }`
- `reactivate` → `{ id, status: "active" }`
- `terminate` → `{ id, status: "terminated" }`

#### `PUT /v1/agents/:id` — UpdateAgentRequest → 200 `{ message }`

Tất cả field `omitempty` — FE **chỉ gửi field user thực sự edit**, không re-submit masked PII:

```ts
interface UpdateAgentRequest {
  full_name?: string;       // min 1, max 255
  phone?: string;           // vn_phone
  email?: string;           // email, max 320
  province?: string;        // min 1, max 100
  district?: string;        // max 100
  business_type?: 'COMPANY' | 'HOUSEHOLD' | 'INDIVIDUAL';
  address_line?: string;    // max 500
  province_code?: string;   // max 10
  district_code?: string;   // max 10
  ward_code?: string;       // max 10
  partner_id?: string;      // max 100
  user_id?: string;         // uuid
  referral_code?: string;   // max 50
  referral_name?: string;   // max 255
  maker_id?: string;        // uuid
  checker_id?: string;      // uuid
}
```

#### `PUT /v1/agents/:id/bank` — UpdateAgentBankRequest → 200 `{ message }`

```ts
interface UpdateAgentBankRequest {
  bank_name: string;            // required, max 100
  bank_account_number: string;  // required, max 50 — sẽ mask "****1234" khi GET
  bank_account_holder: string;  // required, max 255
}
```

#### `POST /v1/agents/:id/documents` — UploadDocumentRequest → 201

```ts
interface UploadDocumentRequest {
  document_type:                // required
    | 'BUSINESS_LICENSE' | 'ID_CARD_FRONT' | 'ID_CARD_BACK'
    | 'SELFIE' | 'TAX_CERTIFICATE' | 'CONTRACT'
    | 'BHXH_AUTH' | 'BANK_STATEMENT' | 'OTHER';
  file_url: string;             // required — pre-signed S3 URL
  file_name?: string;           // max 255
  file_size?: number;           // bytes, 1..10_485_760 (10MB)
  mime_type?: 'image/jpeg' | 'image/png' | 'image/heic' | 'application/pdf';
  checksum?: string;            // max 64
  issued_at?: string;           // ISO 8601
  expires_at?: string;          // ISO 8601
  kyc_submission_id?: string;   // uuid
}
```

**Response:** `{ id, agent_id, document_type, status, created_at }`.

#### `PUT /v1/documents/:id/review` — ReviewDocumentRequest → 200 `{ message }`

```ts
interface ReviewDocumentRequest {
  action: 'approve' | 'reject';  // required
  rejection_reason?: string;     // required khi action='reject'
}
```

#### `GET /v1/agents/:id/documents` → 200 cursor list

Response: `{ data: [{ id, document_type, status, url, uploaded_at, ... }], cursor, has_more }`.

---

### 14.2 Onboarding / KYC module

#### `POST /v1/onboarding/agents` — OnboardingAgentRequest → 200 OnboardingAgentResponse

Tạo MasterAgent (L0) all-in-one (identity + profile + attachments). `parent_id` optional cho L0.

```ts
interface OnboardingAgentRequest {
  // Identity
  full_name: string;               // required, 1–255
  phone: string;                   // required, 9–15 (không dùng vn_phone ở đây)
  email: string;                   // required, email, max 320
  province: string;                // required, 1–100
  district?: string;               // max 100
  parent_id?: string;              // uuid, optional cho L0
  referral_code?: string;          // max 50
  referral_name?: string;          // max 255

  // Profile (optional)
  business_name?: string;          // max 255
  business_type?: string;          // max 100
  tax_code?: string;               // max 20
  representative_name?: string;    // max 255
  representative_cccd?: string;    // max 20, plaintext ở payload, encrypt khi promote
  bank_name?: string;              // max 255
  bank_account?: string;           // max 50, plaintext ở payload, encrypt khi promote
  permanent_address?: string;      // max 500
  contact_address?: string;        // max 500

  // Attachments (S3 URLs)
  business_license_url?: string;   // url
  agency_contract_url?: string;    // url
  ekyc_cccd_front_url?: string;    // url
  ekyc_cccd_back_url?: string;     // url
  ekyc_selfie_url?: string;        // url

  mode?: 'submit' | 'draft';       // default "draft"; "submit" auto-submit nếu đủ docs
}

interface OnboardingAgentResponse {
  agent_id: string;      // pre-allocated UUID
  ticket_id: string;     // "TK-YYYYMMDD-XXXXXXXX"
  status: 'draft' | 'pending_review';
  message?: string;
}
```

#### `POST /v1/onboarding/tickets` → 200 CreateTicketResponse

Request (inline):
```ts
{ agent_id?: string; agent_payload?: AgentPayload }
```

Response:
```ts
interface CreateTicketResponse {
  id: string;          // internal UUID
  ticket_id: string;   // "TK-YYYYMMDD-XXXXXXXX"
  agent_id: string;
  status: 'draft';
}
```

#### `GET /v1/onboarding/tickets` → 200 cursor list của `TicketListItemDTO` (PII masked)

```ts
interface TicketListItemDTO {
  id: string;
  ticket_id: string;
  agent_id: string;
  agent_name: string;
  status: 'draft' | 'pending_review' | 'approved' | 'rejected';
  current_step?: string;
  full_name?: string;
  phone?: string;       // "090****890"
  email?: string;       // "a***@example.com"
  cccd?: string;        // "****8901"
  created_at: string;
  updated_at: string;
}
```

#### `GET /v1/onboarding/tickets/:ticket_id` → 200 TicketDetailResponse

```ts
interface TicketDetailResponse {
  id: string;
  ticket_id: string;
  agent_id: string;
  status: string;
  full_name?: string;
  phone?: string;
  email?: string;
  cccd?: string;
  workflow_id?: string;      // Temporal
  run_id?: string;
  current_step?: string;     // vd "awaiting_review"
  maker_id?: string;
  checker_id?: string;
  reject_reason?: string;
  agent_payload?: string;    // full JSON snapshot (schema v1)
  attachments: AttachmentDTO[];
  submitted_at?: string;
  reviewed_at?: string;
  created_at: string;
  updated_at: string;
}
```

**`AgentPayload` (JSON snapshot, schema_version = "v1"):**

```ts
interface AgentPayload {
  schema_version: 'v1';
  agent_id?: string;
  identity: {
    full_name?: string;
    phone?: string;        // plaintext trong payload
    email?: string;        // plaintext trong payload
    province?: string;
    district?: string;
    parent_id?: string;
    referral_code?: string;
    referral_name?: string;
  };
  profile: {
    business_name?: string;
    business_type?: string;
    tax_code?: string;
    representative_name?: string;
    representative_cccd?: string;
    bank_name?: string;
    bank_account?: string;
    permanent_address?: string;
    contact_address?: string;
  };
  attachments: Array<{
    type: 'business_license' | 'agency_contract' | 'ekyc_cccd_front' | 'ekyc_cccd_back' | 'ekyc_selfie';
    url: string;
    uploaded_at?: string;
  }>;
  mode?: 'submit' | 'draft';
}
```

#### `POST /v1/onboarding/tickets/:ticket_id/attachments` — UploadAttachmentRequest

```ts
interface UploadAttachmentRequest {
  type: 'business_license' | 'agency_contract' | 'ekyc_photos';  // required
  url: string;                                                    // required, url
}
```

#### `PUT /v1/onboarding/tickets/:ticket_id/submit` — SubmitTicketRequest

```ts
interface SubmitTicketRequest { maker_id?: string }  // uuid
```

Response: `{ ticket_id, status: "pending_review", workflow_id, run_id, current_step: "awaiting_review" }`.

#### `PUT /v1/onboarding/tickets/:ticket_id/review` — ReviewTicketRequest

```ts
interface ReviewTicketRequest {
  action: 'approve' | 'reject';  // required
  reason?: string;                // max 1000, required khi action='reject' (handler-level check)
  checker_id?: string;            // uuid
}
```

Response: `{ ticket_id, status: "approved"|"rejected", action }`.

#### `POST /v1/agents/:id/kyc` — **multipart/form-data** → 201

- Fields: `id_number` (form, required), `id_front` (file), `id_back` (file), `selfie` (file).
- Max payload: 32MB in memory.
- Response: `{ submission_id, status: "pending", workflow_id }`.

#### `GET /v1/kyc/pending` → 200 (reviewer-only)

Query: `cursor?`, `page_size?`. Response: array of `{ submission_id, agent_id, tier, status: "pending", submitted_at }`.

#### `GET /v1/kyc/:submission_id` → 200

`{ submission_id, agent_id, tier, status, submitted_at, reviewer_id, reviewer_note }`.

#### `PUT /v1/kyc/:submission_id/review` — reviewKYCRequest

```ts
interface ReviewKYCRequest {
  action: 'approve' | 'reject';   // required
  reviewer_note?: string;
}
```

Response: `{ submission_id, action, status: "signal_sent" }`.

---

### 14.3 Catalog module

#### `GET /v1/product-categories` → 200 array

```ts
interface ProductCategoryDTO {
  id: string;
  name: string;
  slug: string;
  description: string;
  min_kyc_tier: number;       // 0 | 1 | 2
  display_order: number;
  accessible: boolean;        // theo agent hiện tại
}
```

#### `GET /v1/product-categories/:id/products` → 200 cursor list

```ts
interface ProductDTO {
  id: string;
  name: string;
  code: string;
  description: string;
  supplier_code: string;
  reference_price: number;
  is_active: boolean;
}
```

#### `PUT /v1/agents/:id/product-access/:category_id` — toggleAccessRequest

```ts
interface ToggleAccessRequest { action: 'enable' | 'disable' }
```

Response: `{ agent_id, category_id, action }`.

#### `PUT /v1/agents/:id/commission-overrides/:category_id` — setOverrideRequest

```ts
interface SetOverrideRequest {
  custom_rate_pct: number;   // required, 0 ≤ x ≤ 15
}
```

Response: `{ agent_id, category_id, custom_rate_pct }`.

#### `GET /v1/agents/:id/commission-overrides` → 200 array

```ts
interface CommissionOverrideDTO {
  agent_id: string;
  category_id: string;
  rate_pct: number;
}
```

---

### 14.4 Order module

#### `POST /v1/orders` — CreateOrderRequest → 201 CreateOrderResponse

```ts
interface CreateOrderRequest {
  product_id: string;                           // required, uuid
  quantity: number;                             // required, 1–10000
  customer_name: string;                        // required, 1–255 (plaintext, BE encrypt)
  customer_phone: string;                       // required, vn_phone
  payment_method: 'UPC' | 'Cash App';           // required
  fee?: number;                                 // gte=0, default 0
  vat?: number;                                 // gte=0, default 0
  payment_meta?: unknown;                       // arbitrary JSON
  metadata?: unknown;
  order_meta?: unknown;
  agent_id?: string;                            // uuid, BFF inject từ JWT
  user_id?: string;                             // uuid
  product_category_id?: string;                 // uuid
  supplier_id?: string;                         // uuid
  total_amount?: number;                        // gte=0, auto-calc nếu omit
  currency?: 'VND' | 'USD';                     // default "VND"
  order_ref_id?: string;
}

interface CreateOrderResponse {
  order_id: string;
  order_code: string;      // human-readable
  status: string;          // "pending_payment"
  total_amount: number;
  currency: 'VND' | 'USD';
}
```

#### `POST /v1/orders/:id/pay` — PayRequest → 200

```ts
interface PayRequest {
  payment_method: 'vnpay' | 'momo';   // required
}
```

**Response (vnpay):** `{ payment_url: string }` — redirect user.
**Response (momo):** `{ payment_url, deeplink, qr_code_url }`.

#### `POST /v1/orders/:id/refund` — RefundOrderRequest → 200

```ts
interface RefundOrderRequest {
  amount?: number;    // gte=0, 0 = full refund
  reason?: string;    // max 255
}
```

Response: `{ order_id, status: "refund_requested", amount }`.

#### `OrderDTO` (shape tham chiếu cho list/detail khi được expose)

```ts
interface OrderDTO {
  id: string;
  order_code: string;
  agent_id: string;
  user_id?: string;
  product_id: string;
  product_category_id: string;
  supplier_id: string;
  customer_name?: string;       // PII, conditional
  customer_phone?: string;      // PII, conditional
  quantity: number;
  total_amount: number;
  fee: number;
  vat: number;
  currency: 'VND' | 'USD';
  status: 'pending_payment' | 'paid' | 'processing' | 'completed'
        | 'cancelled' | 'refund_pending' | 'refunded';
  payment_method: string;
  payment_time?: string;
  refund_status: 'none' | 'pending' | 'completed' | 'failed';
  refund_time?: string;
  payment_meta?: unknown;
  metadata?: unknown;
  order_meta?: unknown;
  order_ref_id?: string;
  created_at: string;
}
```

#### Webhooks (BE-to-BE, FE không gọi — chỉ document để hiểu shape)

- `POST /v1/webhooks/payment/vnpay` — VNPay IPN, x-www-form-urlencoded, HMAC-SHA512 qua `vnp_SecureHash`. Response `{ "RspCode": "00", "Message": "Confirm Success" }`.
- `POST /v1/webhooks/payment/momo` — MoMo IPN, JSON, HMAC-SHA256 qua `signature`. Response `{ "resultCode": 0, "message": "ok" }`.
- `POST /v1/webhooks/supplier/:code` — raw body + header `X-Event-Id`/`X-Signature`/`X-Timestamp`. Response `{ "result": "ok" }`.

**Dedup keys:** VNPay `vnp_TransactionNo` → `vnp_TxnRef`; MoMo `transId` → `requestId`; Supplier `X-Event-Id` → SHA256(body).

---

### 14.5 Validation cheat-sheet

| Tag | Ý nghĩa | Ví dụ |
|---|---|---|
| `required` | Bắt buộc, non-empty | |
| `omitempty` | Optional — bỏ qua khi rỗng | |
| `uuid` | UUID v4 | `"3fa85f64-..."` |
| `email` | RFC 5322 email | |
| `url` | URL hợp lệ | |
| `min=N` / `max=N` | String length hoặc numeric bound | |
| `gte=N` / `lte=N` | Numeric ≥ / ≤ | `gte=0,lte=15` |
| `oneof=A B C` | Một trong danh sách (space-separated) | `oneof=approve reject` |
| `vn_phone` | `^(03\|05\|07\|08\|09)\d{8}$` | `0901234567` |
| `status_change_reason` | Required (1–500) khi `action ∈ {suspend, terminate}` | |

### 14.6 Quy tắc khi FE submit form

1. **PII mask ≠ plaintext:** không re-submit `"090****890"` trong update form — nếu user không edit phone, **omit field** thay vì gửi lại giá trị đang hiển thị.
2. **UUID path params:** `:id`, `:ticket_id`, `:submission_id`, `:category_id` phải là UUID v4 hợp lệ.
3. **Idempotency:** mọi `POST/PUT/PATCH/DELETE` cần `X-Idempotency-Key` (xem §6). Thiếu → `code 3002`.
4. **Maker/Checker:** BFF (gin-vue-admin) tự inject `maker_id`/`checker_id` từ JWT — FE không phải gửi.
5. **Multipart vs JSON:** chỉ `POST /v1/agents/:id/kyc` dùng `multipart/form-data`. Tất cả document upload khác (`POST /v1/agents/:id/documents`, onboarding attachments) dùng **pre-signed S3 URL** — FE upload thẳng lên S3 rồi gửi URL cho BE.
6. **Webhook URL path param:** đường dẫn `/v1/webhooks/payment/:gateway` — `:gateway` là literal giá trị `vnpay` hoặc `momo`, không phải API Gateway service.

---

*End of guide. Vấn đề chưa cover → file Issue ở repo backend với label `docs:integration`.*
