package system

import (
	"net/http"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SkillsApi struct{}

func (s *SkillsApi) GetTools(c *gin.Context) {
	data, err := skillsService.Tools(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("Failed to get tool list", zap.Error(err))
		response.FailWithMessage("Failed to get tool list", c)
		return
	}
	response.OkWithDetailed(gin.H{"tools": data}, "Retrieved successfully", c)
}

func (s *SkillsApi) GetSkillList(c *gin.Context) {
	var req request.SkillToolRequest
	_ = c.ShouldBindJSON(&req)
	data, err := skillsService.List(c.Request.Context(), req.Tool)
	if err != nil {
		global.GVA_LOG.Error("Failed to get skill list", zap.Error(err))
		response.FailWithMessage("Failed to get skill list", c)
		return
	}
	response.OkWithDetailed(gin.H{"skills": data}, "Retrieved successfully", c)
}

func (s *SkillsApi) GetSkillDetail(c *gin.Context) {
	var req request.SkillDetailRequest
	_ = c.ShouldBindJSON(&req)
	data, err := skillsService.Detail(c.Request.Context(), req.Tool, req.Skill)
	if err != nil {
		global.GVA_LOG.Error("Failed to get skill detail", zap.Error(err))
		response.FailWithMessage("Failed to get skill detail", c)
		return
	}
	response.OkWithDetailed(gin.H{"detail": data}, "Retrieved successfully", c)
}

func (s *SkillsApi) SaveSkill(c *gin.Context) {
	var req request.SkillSaveRequest
	_ = c.ShouldBindJSON(&req)
	if err := skillsService.Save(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to save skill", zap.Error(err))
		response.FailWithMessage("Failed to save skill", c)
		return
	}
	response.OkWithMessage("Saved successfully", c)
}

func (s *SkillsApi) DeleteSkill(c *gin.Context) {
	var req request.SkillDeleteRequest
	_ = c.ShouldBindJSON(&req)
	if err := skillsService.Delete(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to delete skill", zap.Error(err))
		response.FailWithMessage("Failed to delete skill: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

func (s *SkillsApi) CreateScript(c *gin.Context) {
	var req request.SkillScriptCreateRequest
	_ = c.ShouldBindJSON(&req)
	fileName, content, err := skillsService.CreateScript(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to create script", zap.Error(err))
		response.FailWithMessage("Failed to create script", c)
		return
	}
	response.OkWithDetailed(gin.H{"fileName": fileName, "content": content}, "Created successfully", c)
}

func (s *SkillsApi) GetScript(c *gin.Context) {
	var req request.SkillFileRequest
	_ = c.ShouldBindJSON(&req)
	content, err := skillsService.GetScript(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to read script", zap.Error(err))
		response.FailWithMessage("Failed to read script", c)
		return
	}
	response.OkWithDetailed(gin.H{"content": content}, "Retrieved successfully", c)
}

func (s *SkillsApi) SaveScript(c *gin.Context) {
	var req request.SkillFileSaveRequest
	_ = c.ShouldBindJSON(&req)
	if err := skillsService.SaveScript(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to save script", zap.Error(err))
		response.FailWithMessage("Failed to save script", c)
		return
	}
	response.OkWithMessage("Saved successfully", c)
}

func (s *SkillsApi) CreateResource(c *gin.Context) {
	var req request.SkillResourceCreateRequest
	_ = c.ShouldBindJSON(&req)
	fileName, content, err := skillsService.CreateResource(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to create resource", zap.Error(err))
		response.FailWithMessage("Failed to create resource", c)
		return
	}
	response.OkWithDetailed(gin.H{"fileName": fileName, "content": content}, "Created successfully", c)
}

func (s *SkillsApi) GetResource(c *gin.Context) {
	var req request.SkillFileRequest
	_ = c.ShouldBindJSON(&req)
	content, err := skillsService.GetResource(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to read resource", zap.Error(err))
		response.FailWithMessage("Failed to read resource", c)
		return
	}
	response.OkWithDetailed(gin.H{"content": content}, "Retrieved successfully", c)
}

func (s *SkillsApi) SaveResource(c *gin.Context) {
	var req request.SkillFileSaveRequest
	_ = c.ShouldBindJSON(&req)
	if err := skillsService.SaveResource(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to save resource", zap.Error(err))
		response.FailWithMessage("Failed to save resource", c)
		return
	}
	response.OkWithMessage("Saved successfully", c)
}

func (s *SkillsApi) CreateReference(c *gin.Context) {
	var req request.SkillReferenceCreateRequest
	_ = c.ShouldBindJSON(&req)
	fileName, content, err := skillsService.CreateReference(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to create reference", zap.Error(err))
		response.FailWithMessage("Failed to create reference", c)
		return
	}
	response.OkWithDetailed(gin.H{"fileName": fileName, "content": content}, "Created successfully", c)
}

func (s *SkillsApi) GetReference(c *gin.Context) {
	var req request.SkillFileRequest
	_ = c.ShouldBindJSON(&req)
	content, err := skillsService.GetReference(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to read reference", zap.Error(err))
		response.FailWithMessage("Failed to read reference", c)
		return
	}
	response.OkWithDetailed(gin.H{"content": content}, "Retrieved successfully", c)
}

func (s *SkillsApi) SaveReference(c *gin.Context) {
	var req request.SkillFileSaveRequest
	_ = c.ShouldBindJSON(&req)
	if err := skillsService.SaveReference(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to save reference", zap.Error(err))
		response.FailWithMessage("Failed to save reference", c)
		return
	}
	response.OkWithMessage("Saved successfully", c)
}

func (s *SkillsApi) CreateTemplate(c *gin.Context) {
	var req request.SkillTemplateCreateRequest
	_ = c.ShouldBindJSON(&req)
	fileName, content, err := skillsService.CreateTemplate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to create template", zap.Error(err))
		response.FailWithMessage("Failed to create template", c)
		return
	}
	response.OkWithDetailed(gin.H{"fileName": fileName, "content": content}, "Created successfully", c)
}

func (s *SkillsApi) GetTemplate(c *gin.Context) {
	var req request.SkillFileRequest
	_ = c.ShouldBindJSON(&req)
	content, err := skillsService.GetTemplate(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to read template", zap.Error(err))
		response.FailWithMessage("Failed to read template", c)
		return
	}
	response.OkWithDetailed(gin.H{"content": content}, "Retrieved successfully", c)
}

func (s *SkillsApi) SaveTemplate(c *gin.Context) {
	var req request.SkillFileSaveRequest
	_ = c.ShouldBindJSON(&req)
	if err := skillsService.SaveTemplate(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to save template", zap.Error(err))
		response.FailWithMessage("Failed to save template", c)
		return
	}
	response.OkWithMessage("Saved successfully", c)
}

func (s *SkillsApi) GetGlobalConstraint(c *gin.Context) {
	var req request.SkillToolRequest
	_ = c.ShouldBindJSON(&req)
	content, exists, err := skillsService.GetGlobalConstraint(c.Request.Context(), req.Tool)
	if err != nil {
		global.GVA_LOG.Error("Failed to read global constraint", zap.Error(err))
		response.FailWithMessage("Failed to read global constraint", c)
		return
	}
	response.OkWithDetailed(gin.H{"content": content, "exists": exists}, "Retrieved successfully", c)
}

func (s *SkillsApi) SaveGlobalConstraint(c *gin.Context) {
	var req request.SkillGlobalConstraintSaveRequest
	_ = c.ShouldBindJSON(&req)
	if err := skillsService.SaveGlobalConstraint(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to save global constraint", zap.Error(err))
		response.FailWithMessage("Failed to save global constraint", c)
		return
	}
	response.OkWithMessage("Saved successfully", c)
}

func (s *SkillsApi) PackageSkill(c *gin.Context) {
	var req request.SkillPackageRequest
	_ = c.ShouldBindJSON(&req)

	fileName, data, err := skillsService.Package(c.Request.Context(), req)
	if err != nil {
		global.GVA_LOG.Error("Failed to package skill", zap.Error(err))
		response.FailWithMessage("Failed to package skill: "+err.Error(), c)
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Data(http.StatusOK, "application/zip", data)
}

func (s *SkillsApi) DownloadOnlineSkill(c *gin.Context) {
	var req request.DownloadOnlineSkillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("Invalid parameters", c)
		return
	}

	if err := skillsService.DownloadOnlineSkill(c.Request.Context(), req); err != nil {
		global.GVA_LOG.Error("Failed to download online skill", zap.Error(err))
		response.FailWithMessage("Failed to download online skill: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Downloaded successfully", c)
}
