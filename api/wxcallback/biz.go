package wxcallback

import (
	"bytes"
	"encoding/json"
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
	type Message struct {
		ToUserName   string `json:"ToUserName"`
		FromUserName string `json:"FromUserName"`
		CreateTime   int64  `json:"CreateTime"`
		MsgType      string `json:"MsgType"`
		Content      string `json:"Content"`
		MsgId        int64  `json:"MsgId"`
	}
	var msg Message
	err := json.Unmarshal([]byte(r.PostBody), &msg)
	if err != nil {
		log.Error(err)
		return err
	}

	if r.MsgType != "text" {
		return nil
	}
	bot, err := dao.GetTalksAIbot(r.Appid)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("查询到该公众号有绑定AI客服 %+v, %s <-", bot, msg.Content)
	// 查询是否有自动回复的配置，包含是否要求关键字、前缀、后缀
	gptReplyIfNeeded(bot, msg.FromUserName, msg.Content, token)
	return nil
}

func gptReplyIfNeeded(bot *model.TalksAIBot, toUser, question, token string) {
	if bot.Filters != "" {
		//Check filter
	}
	reqGpt := map[string]interface{}{
		"sessionId": toUser,
		"question":  question,
		"botId":     bot.BotID,
		"dec":       true,
	}

	jsonData, err := json.Marshal(reqGpt)
	if err != nil {
		log.Errorf("JSON编码失败: %+v", err)
		return
	}

	url := "https://backend.talks-ai.com/api/chat"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("发送消息到talks ai 失败 step 1 %+v", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("发送消息到talks ai 失败 step 2 %+v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("发送消息到talks ai 失败 step 3 %+v", err)
		return
	}
	var json struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	if err := binding.JSON.BindBody(body, &json); err != nil {
		log.Errorf("发送talks ai失败 %+v", err)
		return
	}
	log.Infof("发送talks ai 结果 %+v", json)
	if json.Code == 0 {
		postContent(toUser, json.Data, token)
	}
}

func postContent(to, content string, token string) {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}
	data["touser"] = to
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("JSON编码失败: %+v", err)
		return
	}

	log.Infof("send content %s to %+v", content, jsonData)

	// 创建POST请求
	url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + token
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("发送消息到公众号失败 step 1 %+v", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("发送消息到公众号失败 step 2 %+v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("发送消息到公众号失败 step 3 %+v", err)
		return
	}

	var json interface{}
	if err := binding.JSON.BindBody(body, &json); err != nil {
		log.Errorf("发送消息到公众号失败 %+v", err)
		return
	}
	log.Infof("发送公众号消息结果 %+v", body)
}
