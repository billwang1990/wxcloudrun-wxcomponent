package talksai

import (
	"github.com/gin-gonic/gin"
)

// Routers 路由
func Routers(e *gin.Engine) {
	g := e.Group("/ai")
	g.POST("/bindingwx/:botid", BindBot)
	g.PUT("/bindingwx/:botid", UpdateBot)
	g.GET("/bindingwx/:botid", QueryBoundBot)
	g.DELETE("/bindingwx/:botid", DeteleBoundBot)
}
