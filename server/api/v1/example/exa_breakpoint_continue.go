package example

import (
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/model/example"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	exampleRes "github.com/huuhoait/gin-vue-admin/server/model/example/response"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BreakpointContinue
// @Tags      ExaFileUploadAndDownload
// @Summary   Breakpoint resume upload to server
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                           true  "an example for breakpoint resume upload"
// @Success   200   {object}  response.Response{msg=string}  "Breakpoint resume upload to server"
// @Router    /fileUploadAndDownload/breakpointContinue [post]
func (b *FileUploadAndDownloadApi) BreakpointContinue(c *gin.Context) {
	fileMd5 := c.Request.FormValue("fileMd5")
	fileName := c.Request.FormValue("fileName")
	chunkMd5 := c.Request.FormValue("chunkMd5")
	chunkNumber, _ := strconv.Atoi(c.Request.FormValue("chunkNumber"))
	chunkTotal, _ := strconv.Atoi(c.Request.FormValue("chunkTotal"))
	_, FileHeader, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("Failed to receive file!", zap.Error(err))
		response.FailWithMessage("Failed to receive file", c)
		return
	}
	f, err := FileHeader.Open()
	if err != nil {
		global.GVA_LOG.Error("Failed to read file!", zap.Error(err))
		response.FailWithMessage("Failed to read file", c)
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	cen, _ := io.ReadAll(f)
	if !utils.CheckMd5(cen, chunkMd5) {
		global.GVA_LOG.Error("Failed to verify md5!", zap.Error(err))
		response.FailWithMessage("MD5 verification failed", c)
		return
	}
	file, err := fileUploadAndDownloadService.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		global.GVA_LOG.Error("Failed to find or create record!", zap.Error(err))
		response.FailWithMessage("Failed to find or create record", c)
		return
	}
	pathC, err := utils.BreakPointContinue(cen, fileName, chunkNumber, chunkTotal, fileMd5)
	if err != nil {
		global.GVA_LOG.Error("Failed to resume upload!", zap.Error(err))
		response.FailWithMessage("Breakpoint resume upload failed", c)
		return
	}

	if err = fileUploadAndDownloadService.CreateFileChunk(file.ID, pathC, chunkNumber); err != nil {
		global.GVA_LOG.Error("Failed to create file record!", zap.Error(err))
		response.FailWithMessage("Failed to create file record", c)
		return
	}
	response.OkWithMessage("Chunk created successfully", c)
}

// FindFile
// @Tags      ExaFileUploadAndDownload
// @Summary   Find file
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                        true  "Find the file"
// @Success   200   {object}  response.Response{data=exampleRes.FileResponse,msg=string}  "Find file, returns file details"
// @Router    /fileUploadAndDownload/findFile [get]
func (b *FileUploadAndDownloadApi) FindFile(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	chunkTotal, _ := strconv.Atoi(c.Query("chunkTotal"))
	file, err := fileUploadAndDownloadService.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		global.GVA_LOG.Error("Failed to find!", zap.Error(err))
		response.FailWithMessage("Search failed", c)
	} else {
		response.OkWithDetailed(exampleRes.FileResponse{File: file}, "Found successfully", c)
	}
}

// BreakpointContinueFinish
// @Tags      ExaFileUploadAndDownload
// @Summary   Create file
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                            true  "File upload complete"
// @Success   200   {object}  response.Response{data=exampleRes.FilePathResponse,msg=string}  "Create file, returns file path"
// @Router    /fileUploadAndDownload/findFile [post]
func (b *FileUploadAndDownloadApi) BreakpointContinueFinish(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	filePath, err := utils.MakeFile(fileName, fileMd5)
	if err != nil {
		global.GVA_LOG.Error("Failed to create file!", zap.Error(err))
		response.FailWithDetailed(exampleRes.FilePathResponse{FilePath: filePath}, "File creation failed", c)
	} else {
		response.OkWithDetailed(exampleRes.FilePathResponse{FilePath: filePath}, "File created successfully", c)
	}
}

// RemoveChunk
// @Tags      ExaFileUploadAndDownload
// @Summary   Remove chunk
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                           true  "Remove cached chunk"
// @Success   200   {object}  response.Response{msg=string}  "Remove chunk"
// @Router    /fileUploadAndDownload/removeChunk [post]
func (b *FileUploadAndDownloadApi) RemoveChunk(c *gin.Context) {
	var file example.ExaFile
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Path traversal prevention
	if strings.Contains(file.FilePath, "..") || strings.Contains(file.FilePath, "../") || strings.Contains(file.FilePath, "./") || strings.Contains(file.FilePath, ".\\") {
		response.FailWithMessage("Illegal path, deletion forbidden", c)
		return
	}
	err = utils.RemoveChunk(file.FileMd5)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete cached chunk!", zap.Error(err))
		return
	}
	err = fileUploadAndDownloadService.DeleteFileChunk(file.FileMd5, file.FilePath)
	if err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Cached chunks deleted successfully", c)
}
