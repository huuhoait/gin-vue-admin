package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type CustomerRouter struct{}

func (e *CustomerRouter) InitCustomerRouter(Router *gin.RouterGroup) {
	customerRouter := Router.Group("customer").Use(middleware.OperationRecord())
	customerRouterWithoutRecord := Router.Group("customer")
	{
		customerRouter.POST("customer", exaCustomerApi.CreateExaCustomer)   // CreateCustomer
		customerRouter.PUT("customer", exaCustomerApi.UpdateExaCustomer)    // updateCustomer
		customerRouter.DELETE("customer", exaCustomerApi.DeleteExaCustomer) // deleteCustomer
	}
	{
		customerRouterWithoutRecord.GET("customer", exaCustomerApi.GetExaCustomer)         // getDocumentOneCustomerInformation
		customerRouterWithoutRecord.GET("customerList", exaCustomerApi.GetExaCustomerList) // getCustomerList
	}
}
