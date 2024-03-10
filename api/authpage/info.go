package authpage

import (
	"encoding/json"
	"net/http"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	wxbase "github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx/base"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/dao"
	"github.com/gin-gonic/gin"
)

func getComponentInfoHandler(c *gin.Context) {
	value := dao.GetCommKv("authinfo", "{}")
	var mapResult map[string]interface{}
	if err := json.Unmarshal([]byte(value), &mapResult); err != nil {
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}
	mapResult["appid"] = wxbase.GetAppid()
	log.Infof("wyq getComponentInfoHandler %+v",mapResult)
	c.JSON(http.StatusOK, errno.OK.WithData(mapResult))
}
