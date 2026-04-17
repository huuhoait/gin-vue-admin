package example

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	common "github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/example"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AttachmentCategoryApi struct{}

// GetCategoryList
// @Tags      GetCategoryList
// @Summary   Media library category list
// @Security  AttachmentCategory
// @Produce   application/json
// @Success   200   {object}  response.Response{data=example.ExaAttachmentCategory,msg=string}  "Media library category list"
// @Router    /attachmentCategory/getCategoryList [get]
func (a *AttachmentCategoryApi) GetCategoryList(c *gin.Context) {
	res, err := attachmentCategoryService.GetCategoryList()
	if err != nil {
		global.GVA_LOG.Error("Failed to get category list!", zap.Error(err))
		response.FailWithMessage("Failed to get category list", c)
		return
	}
	response.OkWithData(res, c)
}

// AddCategory
// @Tags      AddCategory
// @Summary   Add media library category
// @Security  AttachmentCategory
// @accept    application/json
// @Produce   application/json
// @Param     data  body      example.ExaAttachmentCategory  true  "Media library category data"// @Success   200   {object}  response.Response{msg=string}   "Add media library category"
// @Router    /attachmentCategory/addCategory [post]
func (a *AttachmentCategoryApi) AddCategory(c *gin.Context) {
	var req example.ExaAttachmentCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("Invalid parameters!", zap.Error(err))
		response.FailWithMessage("Invalid parameters", c)
		return
	}

	if err := attachmentCategoryService.AddCategory(&req); err != nil {
		global.GVA_LOG.Error("Failed to create/update!", zap.Error(err))
		response.FailWithMessage("Failed to create/update: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Created/updated successfully", c)
}

// DeleteCategory
// @Tags      DeleteCategory
// @Summary   Delete category
// @Security  AttachmentCategory
// @accept    application/json
// @Produce   application/json
// @Param     data  body      common.GetById                true  "Category ID"
// @Success   200   {object}  response.Response{msg=string}  "Delete category"
// @Router    /attachmentCategory/deleteCategory [post]
func (a *AttachmentCategoryApi) DeleteCategory(c *gin.Context) {
	var req common.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Invalid parameters", c)
		return
	}

	if req.ID == 0 {
		response.FailWithMessage("Invalid parameters", c)
		return
	}

	if err := attachmentCategoryService.DeleteCategory(&req.ID); err != nil {
		response.FailWithMessage("Failed to delete", c)
		return
	}

	response.OkWithMessage("Deleted successfully", c)
}
