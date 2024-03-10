package wxcallback

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/dao"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type wxCallbackBizRecord struct {
	CreateTime int64  `json:"CreateTime"`
	ToUserName string `json:"ToUserName"`
	MsgType    string `json:"MsgType"`
	Event      string `json:"Event"`
}

func bizHandler(c *gin.Context) {
	// 记录到数据库
	body, _ := ioutil.ReadAll(c.Request.Body)
	var json wxCallbackBizRecord
	if err := binding.JSON.BindBody(body, &json); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	r := model.WxCallbackBizRecord{
		CreateTime:  time.Unix(json.CreateTime, 0),
		ReceiveTime: time.Now(),
		Appid:       c.Param("appid"),
		ToUserName:  json.ToUserName,
		MsgType:     json.MsgType,
		Event:       json.Event,
		PostBody:    string(body),
	}
	if json.CreateTime == 0 {
		r.CreateTime = time.Unix(1, 0)
	}
	
	log.Infof("bound body %+v record to store is %+v <-------", json, r.PostBody)

	for k, v := range c.Request.Header {
		log.Debugf("xxxxxxxxxxx  %s %s", k, v)
	}
	if err := dao.AddBizCallBackRecord(&r); err != nil {
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}

	token, err := wx.BizGetComponentAccessToken(r.Appid)
	if err == nil {
		go replyMsgIfNeeded(&r, token)
	} else {
		log.Errorf("获取appid: %s 的token失败: %+v", r.Appid, err)
	}
	// 转发到用户配置的地址
	proxyOpen, err := proxyCallbackMsg("", json.MsgType, json.Event, string(body), c)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}
	if !proxyOpen {
		c.String(http.StatusOK, "success")
	}
}

func replyMsgIfNeeded(r *model.WxCallbackBizRecord, token string) error {
	if r.MsgType != "text" {
		return nil
	}
	bot, err := dao.GetTalksAIbot(r.Appid)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("查询到该公众号有绑定AI客服 %+v", bot)

	// 查询是否有自动回复的配置，包含是否要求关键字、前缀、后缀
	// postContent("", token)
	return nil
}

func chatBot() {
	// const data = {
	// 	"sessionId": FromUserName,
	// 	"question": Content,
	// 	"botId": "0705BpLnfgDs",
	// 	"dec": true
	// }
}

func postContent(content string, token string) {

	// 定义要发送的JSON数据
	jsonData := []byte(`{
		"touser": "oDYseuFGkl2rn5zdi_Ve_I6vAwr4",
		"msgtype": "text",
		"text": {
			"content": "\n—————保罗AI客服回复"
		}
	}`)

	// 创建POST请求
	url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + token
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("发送消息到公众失败 step 1 %+v", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("发送消息到公众失败 step 2 %+v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应数据
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("发送消息到公众失败 step 3 %+v", err)
		return
	}
}
