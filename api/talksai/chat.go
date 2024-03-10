package talksai

import (
	"io/ioutil"
	"net/http"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin"
)

type bindingBot struct {
	BotID    int64    `json:"botid"`
	AuthCode string   `json:"code"`
	Filters  []string `json:"filters"`
	Prefix   string   `json:"prefix"`
	Suffix   string   `json:"suffix"`
}

func QueryBindBot(c *gin.Context) {

}

func BindBot(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var json bindingBot
	if err := binding.JSON.BindBody(body, &json); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	log.Infof("Prepare binding bot %+v", json)
	c.JSON(http.StatusOK, errno.OK.WithData(gin.H{"msg": "", "data": "success", "code": 0}))
}
