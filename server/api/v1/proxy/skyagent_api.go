package proxy

import (
	"net/http"

	proxyPkg "github.com/huuhoait/gin-vue-admin/server/service/proxy"
	"github.com/gin-gonic/gin"
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
	doProxy(c, getProxyService().Core.Client, http.MethodGet, "/v1/agents/"+c.Param("id"), nil)
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
	doProxy(c, getProxyService().Core.Client, http.MethodGet, "/v1/admin/agents/"+c.Param("id"), nil)
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
	doProxy(c, getProxyService().Core.Client, http.MethodGet, "/v1/onboarding/tickets/"+c.Param("ticket_id"), nil)
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

func doProxy(c *gin.Context, client *proxyPkg.Client, method, path string, body any) {
	headers := proxyPkg.AuthHeaders(c)

	// Forward query string for GET requests.
	opts := &proxyPkg.RequestOpts{Headers: headers}
	if method == http.MethodGet {
		opts.Query = c.Request.URL.Query()
	}

	envelope, httpStatus, err := client.Do(c.Request.Context(), method, path, body, opts)
	if err != nil {
		proxyPkg.RespondError(c, err)
		return
	}
	proxyPkg.Respond(c, envelope, httpStatus)
}
