package example

import (
	"github.com/gin-gonic/gin"
)

type AttachmentCategoryRouter struct{}

func (r *AttachmentCategoryRouter) InitAttachmentCategoryRouterRouter(Router *gin.RouterGroup) {
	router := Router.Group("attachmentCategory")
	{
		router.GET("getCategoryList", attachmentCategoryApi.GetCategoryList) // category list
		router.POST("addCategory", attachmentCategoryApi.AddCategory)        // Add/editPartClass
		router.POST("deleteCategory", attachmentCategoryApi.DeleteCategory)  // deletePartClass
	}
}
