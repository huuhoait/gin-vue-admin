package plugin

import (
	"github.com/gin-gonic/gin"
)

const (
	OnlyFuncName = "Plugin"
)

// Plugin plugin mode API
type Plugin interface {
	// Register register routes
	Register(group *gin.RouterGroup)

	// RouterPath UserReturnregister routes
	RouterPath() string
}
