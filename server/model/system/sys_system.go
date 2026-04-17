package system

import (
	"github.com/huuhoait/gin-vue-admin/server/config"
)

// configurationFileStructureBody
type System struct {
	Config config.Server `json:"config"`
}
