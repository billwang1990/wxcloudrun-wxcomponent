package talksai

import (
	"io/ioutil"
	"net/http"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/dao"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type bindingBot struct {
	AuthCode string `json:"code"`
	Filters  string `json:"filters"`
	Prefix   string `json:"prefix"`
	Suffix   string `json:"suffix"`
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

	if err := dao.UpdateTalksAIBot(&model.TalksAIBot{
		BotID:   botid,
		Filters: json.Filters,
		Prefix:  json.Prefix,
		Suffix:  json.Suffix,
	}); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	c.JSON(http.StatusOK, errno.OK)
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

	if err := dao.CreateOrUpdateTalksAIBot(&model.TalksAIBot{
		BotID:   botid,
		AppID:   record.AuthorizerAppid,
		Filters: json.Filters,
		Prefix:  json.Prefix,
		Suffix:  json.Suffix,
	}); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}

	c.JSON(http.StatusOK, errno.OK)
	go dao.DeleteCachedAuthorizerAppRecord(record.AuthorizerAppid)
}
