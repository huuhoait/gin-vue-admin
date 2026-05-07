package proxy

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	proxyPkg "github.com/huuhoait/gin-vue-admin/server/service/proxy"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SkyAgentApi provides HTTP handlers that proxy requests to Core/Order services.
type SkyAgentApi struct{}

// ---------------------------------------------------------------------------
// Agents
// ---------------------------------------------------------------------------

// GetAgentList proxies GET /v1/agents with query params.
// @Tags     SkyAgent-Agent
// @Summary  List agents
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    page      query int    false "Page number"
// @Param    pageSize  query int    false "Page size"
// @Param    status    query string false "Filter by status"
// @Param    level     query string false "Filter by level"
// @Param    keyword   query string false "Search keyword"
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/agents [get]
func (s *SkyAgentApi) GetAgentList(c *gin.Context) {
	doProxy(c, getProxyService().Core.Client, http.MethodGet, "/v1/agents", nil)
}

// GetAgentDetail proxies GET /v1/agents/:id.
// @Tags     SkyAgent-Agent
// @Summary  Get agent detail
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    id path string true "Agent ID"
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/agents/{id} [get]
func (s *SkyAgentApi) GetAgentDetail(c *gin.Context) {
	doProxyWith(c, getProxyService().Core.Client, http.MethodGet, "/v1/agents/"+c.Param("id"), nil,
		enrichUserNames("created_by", "updated_by", "maker_id", "checker_id"),
	)
}

// CreateAgent proxies POST /v1/agents.
// @Tags     SkyAgent-Agent
// @Summary  Create a new agent
// @Security ApiKeyAuth
// @Accept   application/json
// @Produce  application/json
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/agents [post]
func (s *SkyAgentApi) CreateAgent(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	doProxy(c, getProxyService().Core.Client, http.MethodPost, "/v1/agents", body)
}

// UpdateAgent proxies PUT /v1/agents/:id.
// @Tags     SkyAgent-Agent
// @Summary  Update agent
// @Security ApiKeyAuth
// @Accept   application/json
// @Produce  application/json
// @Param    id path string true "Agent ID"
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/agents/{id} [put]
func (s *SkyAgentApi) UpdateAgent(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	doProxy(c, getProxyService().Core.Client, http.MethodPut, "/v1/agents/"+c.Param("id"), body)
}

// UpdateAgentStatus proxies PUT /v1/agents/:id/status (approve/suspend/terminate).
// @Tags     SkyAgent-Agent
// @Summary  Update agent status
// @Security ApiKeyAuth
// @Accept   application/json
// @Produce  application/json
// @Param    id path string true "Agent ID"
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/agents/{id}/status [put]
func (s *SkyAgentApi) UpdateAgentStatus(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	doProxy(c, getProxyService().Core.Client, http.MethodPut, "/v1/agents/"+c.Param("id")+"/status", body)
}

// GetAgentAdminDetail proxies GET /v1/admin/agents/:id (full PII — admin only).
// @Tags     SkyAgent-Agent
// @Summary  Get agent admin detail with PII
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    id path string true "Agent ID"
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/agents/{id}/full [get]
func (s *SkyAgentApi) GetAgentAdminDetail(c *gin.Context) {
	doProxyWith(c, getProxyService().Core.Client, http.MethodGet, "/v1/admin/agents/"+c.Param("id"), nil,
		enrichUserNames("created_by", "updated_by", "maker_id", "checker_id"),
	)
}

// ---------------------------------------------------------------------------
// Orders
// ---------------------------------------------------------------------------

// GetOrderList proxies GET /v1/orders with query params.
// @Tags     SkyAgent-Order
// @Summary  List orders
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/orders [get]
func (s *SkyAgentApi) GetOrderList(c *gin.Context) {
	doProxy(c, getProxyService().Order.Client, http.MethodGet, "/v1/orders", nil)
}

// GetOrderDetail proxies GET /v1/orders/:id.
// @Tags     SkyAgent-Order
// @Summary  Get order detail
// @Security ApiKeyAuth
// @Produce  application/json
// @Param    id path string true "Order ID"
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/orders/{id} [get]
func (s *SkyAgentApi) GetOrderDetail(c *gin.Context) {
	doProxy(c, getProxyService().Order.Client, http.MethodGet, "/v1/orders/"+c.Param("id"), nil)
}

// ---------------------------------------------------------------------------
// Products / Suppliers
// ---------------------------------------------------------------------------

// GetProductList proxies GET /v1/products.
// @Tags     SkyAgent-Catalog
// @Summary  List products
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/products [get]
func (s *SkyAgentApi) GetProductList(c *gin.Context) {
	doProxy(c, getProxyService().Core.Client, http.MethodGet, "/v1/products", nil)
}

// GetSupplierList proxies GET /v1/suppliers.
// @Tags     SkyAgent-Catalog
// @Summary  List suppliers
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} proxyPkg.GVAEnvelope
// @Router   /admin-api/v1/suppliers [get]
func (s *SkyAgentApi) GetSupplierList(c *gin.Context) {
	doProxy(c, getProxyService().Core.Client, http.MethodGet, "/v1/suppliers", nil)
}

// ---------------------------------------------------------------------------
// Onboarding Tickets
// ---------------------------------------------------------------------------

func (s *SkyAgentApi) CreateTicket(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	doProxy(c, getProxyService().Core.Client, http.MethodPost, "/v1/onboarding/tickets", body)
}

func (s *SkyAgentApi) ListTickets(c *gin.Context) {
	doProxy(c, getProxyService().Core.Client, http.MethodGet, "/v1/onboarding/tickets", nil)
}

func (s *SkyAgentApi) GetTicket(c *gin.Context) {
	doProxyWith(c, getProxyService().Core.Client, http.MethodGet, "/v1/onboarding/tickets/"+c.Param("ticket_id"), nil,
		enrichUserNames("maker_id", "checker_id"),
	)
}

func (s *SkyAgentApi) UploadTicketAttachment(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	doProxy(c, getProxyService().Core.Client, http.MethodPost, "/v1/onboarding/tickets/"+c.Param("ticket_id")+"/attachments", body)
}

func (s *SkyAgentApi) SubmitTicket(c *gin.Context) {
	doProxy(c, getProxyService().Core.Client, http.MethodPut, "/v1/onboarding/tickets/"+c.Param("ticket_id")+"/submit", nil)
}

func (s *SkyAgentApi) ReviewTicket(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	doProxy(c, getProxyService().Core.Client, http.MethodPut, "/v1/onboarding/tickets/"+c.Param("ticket_id")+"/review", body)
}

// OnboardingAgent proxies POST /v1/onboarding/agents (all-in-one Agent L0 creation).
func (s *SkyAgentApi) OnboardingAgent(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	doProxy(c, getProxyService().Core.Client, http.MethodPost, "/v1/onboarding/agents", body)
}

// ---------------------------------------------------------------------------
// Shared proxy helper
// ---------------------------------------------------------------------------

// injectMakerChecker ensures maker_id/checker_id are set on the forwarded
// body. A nil body becomes {maker_id, checker_id} so promote/review
// endpoints that require checker_id succeed even when the FE sent no body
// (PUT /.../submit, PUT /.../promote). If the FE already provided either
// field, its value is preserved.
func injectMakerChecker(body any, userUUID string) any {
	if userUUID == "" {
		return body
	}
	if body == nil {
		return map[string]any{"maker_id": userUUID, "checker_id": userUUID}
	}
	m, ok := body.(map[string]any)
	if !ok {
		return body
	}
	if _, exists := m["maker_id"]; !exists {
		m["maker_id"] = userUUID
	}
	if _, exists := m["checker_id"]; !exists {
		m["checker_id"] = userUUID
	}
	return m
}

// maxLoggedBodyBytes caps request/response body logging to avoid runaway log
// volume on large payloads (attachment manifests, bulk lists).
const maxLoggedBodyBytes = 8 * 1024

// envelopeTransformer can mutate an envelope after a successful proxy call
// (e.g. to enrich the response data). Must not mutate on error paths.
type envelopeTransformer func(c *gin.Context, env *proxyPkg.GVAEnvelope)

// doProxy is the default proxy path with no response enrichment.
func doProxy(c *gin.Context, client *proxyPkg.Client, method, path string, body any) {
	doProxyWith(c, client, method, path, body, nil)
}

// WARNING: these logs contain full upstream payloads including PII
// (CCCD, phone, bank account, email). They should only be enabled in
// trusted environments. Route admins/ops to the audit chain for
// sanitized, retention-controlled traces instead.
func doProxyWith(c *gin.Context, client *proxyPkg.Client, method, path string, body any, transform envelopeTransformer) {
	headers, err := proxyPkg.AuthHeaders(c)
	if err != nil {
		// JWT middleware should have short-circuited before we got here;
		// if we're still missing a UUID, the auth chain is misconfigured.
		// Fail loudly rather than forwarding an anonymous request.
		global.GVA_LOG.Error("skyagent proxy auth failure",
			zap.String("request_id", middleware.GetRequestID(c)),
			zap.String("method", method),
			zap.String("inbound_path", c.FullPath()),
			zap.String("downstream_path", path),
			zap.Error(err),
		)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": response.ERROR,
			"data": nil,
			"msg":  "unauthorized: missing user identity",
		})
		return
	}

	// Forward query string for GET requests.
	opts := &proxyPkg.RequestOpts{Headers: headers}
	if method == http.MethodGet {
		opts.Query = c.Request.URL.Query()
	}

	// Core validates maker_id / checker_id as uuid in most mutation bodies
	// (CreateAgent, ReviewTicket, promote, ...). BFF is contractually the
	// injection point (see external-frontend-integration.md §14.6 rule 4),
	// so fill them in from JWT when the FE didn't supply a value.
	if method != http.MethodGet {
		body = injectMakerChecker(body, headers["X-Maker-ID"])
	}

	reqID := middleware.GetRequestID(c)
	makerID := headers["X-Maker-ID"]
	start := time.Now()

	startFields := []zap.Field{
		zap.String("request_id", reqID),
		zap.String("maker_id", makerID),
		zap.String("method", method),
		zap.String("inbound_path", c.FullPath()),
		zap.String("downstream_path", path),
	}
	if method == http.MethodGet && opts.Query != nil {
		startFields = append(startFields, zap.String("query", redactQuery(opts.Query)))
	}
	if body != nil {
		startFields = append(startFields, zap.String("request_body", truncatedJSON(body)))
	}
	global.GVA_LOG.Info("skyagent proxy call start", startFields...)

	envelope, httpStatus, err := client.Do(c.Request.Context(), method, path, body, opts)
	durMs := time.Since(start).Milliseconds()

	if err != nil {
		global.GVA_LOG.Error("skyagent proxy call failed",
			zap.String("request_id", reqID),
			zap.String("maker_id", makerID),
			zap.String("method", method),
			zap.String("downstream_path", path),
			zap.Int64("duration_ms", durMs),
			zap.Error(err),
		)
		proxyPkg.RespondError(c, err)
		return
	}

	logFields := []zap.Field{
		zap.String("request_id", reqID),
		zap.String("maker_id", makerID),
		zap.String("method", method),
		zap.String("downstream_path", path),
		zap.Int("http_status", httpStatus),
		zap.Int64("duration_ms", durMs),
	}
	if envelope != nil {
		logFields = append(logFields,
			zap.Int("envelope_code", envelope.Code),
			zap.String("envelope_msg", envelope.Msg),
			zap.String("response_data", truncateBytes(envelope.Data)),
		)
	}
	if httpStatus >= http.StatusBadRequest || (envelope != nil && envelope.Code != 0) {
		global.GVA_LOG.Warn("skyagent proxy call non-ok", logFields...)
	} else {
		global.GVA_LOG.Info("skyagent proxy call done", logFields...)
	}

	// Post-process successful envelopes (e.g. resolve maker/checker UUIDs
	// to display names). Errors inside transformers are not fatal — the
	// original envelope is forwarded regardless.
	if transform != nil && envelope != nil && envelope.Code == 0 && httpStatus < http.StatusBadRequest {
		transform(c, envelope)
	}

	proxyPkg.Respond(c, envelope, httpStatus)
}

// enrichUserNames returns a transformer that resolves any of the given field
// names inside the envelope's top-level data object from UUID -> display name
// (from sys_users) and writes them back as `<field>_name` sidecars. Fields
// that are missing, non-string, or fail to resolve are left untouched.
func enrichUserNames(fields ...string) envelopeTransformer {
	if len(fields) == 0 {
		return nil
	}
	return func(c *gin.Context, env *proxyPkg.GVAEnvelope) {
		if len(env.Data) == 0 || string(env.Data) == "null" {
			return
		}
		var payload map[string]any
		if err := json.Unmarshal(env.Data, &payload); err != nil {
			// Envelope data isn't a JSON object (list/scalar). Skip.
			return
		}

		uuids := make([]string, 0, len(fields))
		for _, f := range fields {
			if v, ok := payload[f].(string); ok && v != "" {
				uuids = append(uuids, v)
			}
		}
		if len(uuids) == 0 {
			return
		}

		names := proxyPkg.ResolveUserNames(c.Request.Context(), uuids)
		if len(names) == 0 {
			return
		}
		for _, f := range fields {
			v, ok := payload[f].(string)
			if !ok || v == "" {
				continue
			}
			if name, found := names[v]; found {
				payload[f+"_name"] = name
			}
		}

		raw, err := json.Marshal(payload)
		if err != nil {
			return
		}
		env.Data = raw
	}
}

// truncatedJSON marshals v and trims to maxLoggedBodyBytes for safe logging.
// Any marshal failure falls back to a fixed sentinel so a logging path never
// panics or aborts the proxy request.
func truncatedJSON(v any) string {
	raw, err := json.Marshal(v)
	if err != nil {
		return "<marshal-error>"
	}
	return truncateBytes(raw)
}

func truncateBytes(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	if len(raw) <= maxLoggedBodyBytes {
		return string(raw)
	}
	return string(raw[:maxLoggedBodyBytes]) + "...<truncated>"
}

// redactQuery renders a query string with obvious secret-bearing params
// (token, key, secret, password) replaced by "***".
func redactQuery(q url.Values) string {
	if len(q) == 0 {
		return ""
	}
	safe := url.Values{}
	for k, vs := range q {
		if isSensitiveParam(k) {
			safe[k] = []string{"***"}
			continue
		}
		safe[k] = vs
	}
	return safe.Encode()
}

func isSensitiveParam(name string) bool {
	switch name {
	case "token", "access_token", "refresh_token", "password", "secret", "api_key", "apikey":
		return true
	}
	return false
}
