// AssumeThisYesinitializeLogicofOnePart

package initialize

import (
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

// initializeGlobalFunctionNumber
func SetupHandlers() {
	// register system reload handler
	utils.GlobalSystemEvents.RegisterReloadHandler(func() error {
		return Reload()
	})
}
