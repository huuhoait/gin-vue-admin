package plugin

import (
	"github.com/gin-gonic/gin"
)

// Plugin plugin mode APIv2
type Plugin interface {
	// Register register routes
	Register(group *gin.Engine)
}
