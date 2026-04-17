package example

import "github.com/huuhoait/gin-vue-admin/server/service"

type ApiGroup struct {
	CustomerApi

	AttachmentCategoryApi
	FileUploadAndDownloadApi
}

var (
	customerService = service.ServiceGroupApp.ExampleServiceGroup.CustomerService

	attachmentCategoryService    = service.ServiceGroupApp.ExampleServiceGroup.AttachmentCategoryService
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
)
