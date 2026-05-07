package router

import (
	"github.com/huuhoait/gin-vue-admin/server/plugin/skyagent/api"

	"github.com/gin-gonic/gin"
)

type SkyAgentRouter struct{}

// Init mounts every SkyAgent BFF route under the supplied private group. The
// group is expected to already carry the JWT + Casbin middleware that gates
// the rest of the admin surface — the plugin's initialize.Router applies
// them. The mount prefix expected by the FE is /admin-api/v1.
func (s *SkyAgentRouter) Init(private *gin.RouterGroup) {
	skyApi := api.Api.SkyAgent
	dashApi := api.Api.Dashboard

	agents := private.Group("agents")
	{
		agents.GET("", skyApi.GetAgentList)
		agents.GET(":id", skyApi.GetAgentDetail)
		agents.POST("", skyApi.CreateAgent)
		agents.PUT(":id", skyApi.UpdateAgent)
		agents.PUT(":id/status", skyApi.UpdateAgentStatus)
		agents.GET(":id/full", skyApi.GetAgentAdminDetail)
	}

	orders := private.Group("orders")
	{
		orders.GET("", skyApi.GetOrderList)
		orders.GET(":id", skyApi.GetOrderDetail)
	}

	products := private.Group("products")
	{
		products.GET("", skyApi.GetProductList)
	}

	suppliers := private.Group("suppliers")
	{
		suppliers.GET("", skyApi.GetSupplierList)
	}

	tickets := private.Group("onboarding/tickets")
	{
		tickets.POST("", skyApi.CreateTicket)
		tickets.GET("", skyApi.ListTickets)
		tickets.GET(":ticket_id", skyApi.GetTicket)
		tickets.POST(":ticket_id/attachments", skyApi.UploadTicketAttachment)
		tickets.PUT(":ticket_id/submit", skyApi.SubmitTicket)
		tickets.PUT(":ticket_id/review", skyApi.ReviewTicket)
	}

	// All-in-one Agent L0 onboarding (Story 11.8)
	private.POST("onboarding/agents", skyApi.OnboardingAgent)

	dash := private.Group("dashboard")
	{
		dash.GET("overview", dashApi.GetDashboardOverview)
	}
}
