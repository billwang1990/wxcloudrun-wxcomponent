package talksai

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/dao"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type bindingBot struct {
	AuthCode       string `json:"code"`
	Filters        string `json:"filters"`
	ExcludeFilters string `json:"excludefilters"`
	Prefix         string `json:"prefix"`
	Suffix         string `json:"suffix"`
}

func QueryBoundBot(c *gin.Context) {
	botid := c.Param("botid")
	bot, err := dao.GetTalksAIbotByBot(botid)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, errno.ErrRequestErr.WithData(err.Error()))
		return
	}
	c.JSON(http.StatusOK, errno.OK.WithData(bot))
}

func DeteleBoundBot(c *gin.Context) {
	botid := c.Param("botid")
	err := dao.DeleteTalksAIBotByBot(botid)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, errno.ErrSystemError)
		return
	}
	c.JSON(http.StatusOK, errno.OK)
}

func UpdateBot(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var json bindingBot
	if err := binding.JSON.BindBody(body, &json); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	botid := c.Param("botid")
	if botid == "" {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData("invalid botid"))
		return
	}
	b := &model.TalksAIBot{
		BotID:          botid,
		Filters:        strings.ReplaceAll(json.Filters, "；", ";"),
		Prefix:         json.Prefix,
		Suffix:         json.Suffix,
		ExcludeFilters: strings.ReplaceAll(json.ExcludeFilters, "；", ";"),
	}
	if err := dao.UpdateTalksAIBot(b); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, errno.OK.WithData(b))
}

func BindBot(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var json bindingBot
	if err := binding.JSON.BindBody(body, &json); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	log.Infof("Prepare binding bot %+v", json)
	botid := c.Param("botid")
	if botid == "" {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData("invalid botid"))
		return
	}
	// query cached appid
	record, err := dao.GetCachedAuthorizerAppRecordByCode(json.AuthCode)
	log.Infof("query cached appid %+v", record)
	if err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}

	defer func () {
		go dao.DeleteCachedAuthorizerAppRecord(record.AuthorizerAppid)
	}()

	if record.VerifyInfo == -1 {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData("该公众号未认证，不能绑定"))
		return
	}
	b := &model.TalksAIBot{
		BotID:   botid,
		AppID:   record.AuthorizerAppid,
		Filters: json.Filters,
		Prefix:  json.Prefix,
		Suffix:  json.Suffix,
		Name:    record.AuthorizerAppName,
		ExcludeFilters: json.ExcludeFilters,
	}
	if err := dao.CreateOrUpdateTalksAIBot(b); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, errno.OK.WithData(b))
	
}
