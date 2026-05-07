package proxy

import (
	v1 "github.com/huuhoait/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type SkyAgentRouter struct{}

// InitSkyAgentRouter registers all SkyAgent proxy routes under a private
// (JWT + Casbin protected) router group.
//
// Route prefix is expected to be /admin-api/v1 (or whatever the caller mounts).
func (s *SkyAgentRouter) InitSkyAgentRouter(Router *gin.RouterGroup) {
	skyApi := v1.ApiGroupApp.ProxyApiGroup

	agents := Router.Group("agents")
	{
		agents.GET("", skyApi.GetAgentList)
		agents.GET(":id", skyApi.GetAgentDetail)
		agents.POST("", skyApi.CreateAgent)
		agents.PUT(":id", skyApi.UpdateAgent)
		agents.PUT(":id/status", skyApi.UpdateAgentStatus)
		agents.GET(":id/full", skyApi.GetAgentAdminDetail)
	}

	orders := Router.Group("orders")
	{
		orders.GET("", skyApi.GetOrderList)
		orders.GET(":id", skyApi.GetOrderDetail)
	}

	products := Router.Group("products")
	{
		products.GET("", skyApi.GetProductList)
	}

	suppliers := Router.Group("suppliers")
	{
		suppliers.GET("", skyApi.GetSupplierList)
	}

	tickets := Router.Group("onboarding/tickets")
	{
		tickets.POST("", skyApi.CreateTicket)
		tickets.GET("", skyApi.ListTickets)
		tickets.GET(":ticket_id", skyApi.GetTicket)
		tickets.POST(":ticket_id/attachments", skyApi.UploadTicketAttachment)
		tickets.PUT(":ticket_id/submit", skyApi.SubmitTicket)
		tickets.PUT(":ticket_id/review", skyApi.ReviewTicket)
	}

	// All-in-one Agent L0 onboarding (Story 11.8)
	Router.POST("onboarding/agents", skyApi.OnboardingAgent)

	dash := Router.Group("dashboard")
	{
		dash.GET("overview", skyApi.GetDashboardOverview)
	}
}
