package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	common "github.com/huuhoait/gin-vue-admin/server/model/common"
)

type AIWorkflowMessage struct {
	ID             string         `json:"id"`
	Role           string         `json:"role"`
	Content        string         `json:"content"`
	Snapshot       common.JSONMap `json:"snapshot"`
	ConversationID string         `json:"conversationId"`
	MessageID      string         `json:"messageId"`
	CreatedAt      string         `json:"createdAt"`
}

type SysAIWorkflowSession struct {
	global.GVA_MODEL
	UserID         uint                `json:"userId" gorm:"column:user_id;index;comment:User ID"`
	Tab            string              `json:"tab" gorm:"column:tab;size:32;index;comment:Sessiontype"`
	Title          string              `json:"title" gorm:"column:title;size:255;comment:SessionTitle"`
	Summary        string              `json:"summary" gorm:"column:summary;type:text;comment:Summary"`
	ConversationID string              `json:"conversationId" gorm:"column:conversation_id;size:255;comment:DifySessionID"`
	MessageID      string              `json:"messageId" gorm:"column:message_id;size:255;comment:DifyMessageID"`
	CurrentNodeID  string              `json:"currentNodeId" gorm:"column:current_node_id;size:64;comment:CurrentChooseInNodeID"`
	Settings       common.JSONMap      `json:"settings" gorm:"column:settings;type:text;comment:Pageset"`
	FormData       common.JSONMap      `json:"formData" gorm:"column:form_data;type:text;comment:TableDocumentData"`
	ResultData     common.JSONMap      `json:"resultData" gorm:"column:result_data;type:text;comment:CurrentDisplayResult"`
	Messages       []AIWorkflowMessage `json:"messages" gorm:"column:messages;serializer:json;type:text;comment:SessionMessage"`
}

func (s *SysAIWorkflowSession) TableName() string {
	return "sys_ai_workflow_sessions"
}
