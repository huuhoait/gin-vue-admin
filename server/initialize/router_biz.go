package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/huuhoait/gin-vue-admin/server/router"
)

// Placeholdermethod, EnsureFileCanByCorrectLoad, AvoidgoEmptyVariableDetectThrow Error, PleaseDo Notdelete.
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]

	holder(publicGroup, privateGroup)

}
