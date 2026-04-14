package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
)

// configurationFileStructureBody
type System struct {
	Config config.Server `json:"config"`
}
