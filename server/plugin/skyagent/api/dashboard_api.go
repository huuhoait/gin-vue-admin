package proxy

import (
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

// DashboardApi provides handlers for the admin dashboard.
type DashboardApi struct{}

// GetDashboardOverview returns aggregated business metrics (cache-aside).
// @Tags     SkyAgent-Dashboard
// @Summary  Get dashboard overview metrics
// @Security ApiKeyAuth
// @Produce  application/json
// @Success  200 {object} response.Response{data=dashboard.OverviewMetrics}
// @Router   /admin-api/v1/dashboard/overview [get]
func (d *DashboardApi) GetDashboardOverview(c *gin.Context) {
	svc := getDashboardService()
	if svc == nil {
		response.FailWithMessage("Dashboard service unavailable", c)
		return
	}
	metrics, err := svc.GetOverview(c.Request.Context())
	if err != nil {
		response.FailWithMessage("Failed to load dashboard", c)
		return
	}
	response.OkWithData(metrics, c)
}
