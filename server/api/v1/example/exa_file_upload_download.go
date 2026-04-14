package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example/request"
	exampleRes "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type FileUploadAndDownloadApi struct{}

// UploadFile
// @Tags      ExaFileUploadAndDownload
// @Summary   Upload file example
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                           true  "Upload file example"
// @Success   200   {object}  response.Response{data=exampleRes.ExaFileResponse,msg=string}  "Upload file example, returns file details"
// @Router    /fileUploadAndDownload/upload [post]
func (b *FileUploadAndDownloadApi) UploadFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	noSave := c.DefaultQuery("noSave", "0")
	_, header, err := c.Request.FormFile("file")
	classId, _ := strconv.Atoi(c.DefaultPostForm("classId", "0"))
	if err != nil {
		global.GVA_LOG.Error("Failed to receive file!", zap.Error(err))
		response.FailWithMessage("Failed to receive file", c)
		return
	}
	file, err = fileUploadAndDownloadService.UploadFile(header, noSave, classId) // Get file path after upload
	if err != nil {
		global.GVA_LOG.Error("Failed to upload file!", zap.Error(err))
		response.FailWithMessage("File upload failed", c)
		return
	}
	response.OkWithDetailed(exampleRes.ExaFileResponse{File: file}, "Uploaded successfully", c)
}

// EditFileName Edit file name or remarks
func (b *FileUploadAndDownloadApi) EditFileName(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = fileUploadAndDownloadService.EditFileName(file)
	if err != nil {
		global.GVA_LOG.Error("Failed to edit!", zap.Error(err))
		response.FailWithMessage("Edit failed", c)
		return
	}
	response.OkWithMessage("Edited successfully", c)
}

// DeleteFile
// @Tags      ExaFileUploadAndDownload
// @Summary   Delete file
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      example.ExaFileUploadAndDownload  true  "Pass in the file ID"
// @Success   200   {object}  response.Response{msg=string}     "Delete file"
// @Router    /fileUploadAndDownload/deleteFile [post]
func (b *FileUploadAndDownloadApi) DeleteFile(c *gin.Context) {
	var file example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := fileUploadAndDownloadService.DeleteFile(file); err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// GetFileList
// @Tags      ExaFileUploadAndDownload
// @Summary   Paginated file list
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.ExaAttachmentCategorySearch                                        true  "Page number, page size, category ID"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Paginated file list, returns list, total, page number, page size"
// @Router    /fileUploadAndDownload/getFileList [post]
func (b *FileUploadAndDownloadApi) GetFileList(c *gin.Context) {
	var pageInfo request.ExaAttachmentCategorySearch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := fileUploadAndDownloadService.GetFileRecordInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Retrieval failed", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Retrieved successfully", c)
}

// ImportURL
// @Tags      ExaFileUploadAndDownload
// @Summary   Import URL
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      example.ExaFileUploadAndDownload  true  "Object"
// @Success   200   {object}  response.Response{msg=string}     "Import URL"
// @Router    /fileUploadAndDownload/importURL [post]
func (b *FileUploadAndDownloadApi) ImportURL(c *gin.Context) {
	var file []example.ExaFileUploadAndDownload
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := fileUploadAndDownloadService.ImportURL(&file); err != nil {
		global.GVA_LOG.Error("Failed to import URL!", zap.Error(err))
		response.FailWithMessage("URL import failed", c)
		return
	}
	response.OkWithMessage("URL imported successfully", c)
}
