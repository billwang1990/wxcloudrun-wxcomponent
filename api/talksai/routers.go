package talksai

import (
	"github.com/gin-gonic/gin"
)

// Routers 路由
func Routers(e *gin.Engine) {
	g := e.Group("/ai")
	g.POST("/binding/:botid", ConfigAutoReply)
}
