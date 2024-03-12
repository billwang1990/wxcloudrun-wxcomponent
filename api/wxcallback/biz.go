package wxcallback

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/dao"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"encoding/xml"
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

	// if err := dao.AddBizCallBackRecord(&r); err != nil {
	// 	c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
	// 	return
	// }

	token, err := wx.BizGetComponentAccessToken(r.Appid)
	if err == nil {
		replyMsgIfNeeded(&r, token, c)
		return
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

func replyMsgIfNeeded(r *model.WxCallbackBizRecord, token string, c *gin.Context) error {
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
	//oDYseuFGkl2rn5zdi_Ve_I6vAwr4 是保罗的
	if msg.FromUserName == "opnbu552g7sy8s63dgm-M60lg7Og" || msg.FromUserName == "oDYseuFGkl2rn5zdi_Ve_I6vAwr4" {
		log.Infof("_____测试被动回复消息1111")
		replyMsg := &ReplyMessage{
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:      "你好我收到了你的消息",
		}

		// // 将回复消息编码为XML格式
		// output, err := xml.MarshalIndent(replyMsg, "", "  ")
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode XML"})
		// 	return nil
		// }

		// 设置响应头并返回XML数据
		log.Infof("_____测试被动回复消息")

		msg, err := xml.Marshal(&replyMsg)
		if err != nil {
			log.Infof("[消息回复] - 将对象进行XML编码出错: %v\n", err)
			return err
		}
		_, _ = c.Writer.Write(msg)
		return nil
	}

	log.Infof("查询到该公众号有绑定AI客服 %+v, %s <-", bot, msg.Content)
	// 查询是否有自动回复的配置，包含是否要求关键字、前缀、后缀
	gptReplyIfNeeded(bot, msg.FromUserName, msg.Content, token)
	return nil
}

func gptReplyIfNeeded(bot *model.TalksAIBot, toUser, question, token string) {
	if bot.Filters != "" {
		//Check filter
		skip := true
		for _, filter := range strings.Split(bot.Filters, ";") {
			if strings.Contains(question, filter) {
				skip = false
				break
			}
		}
		if skip {
			return
		}
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
		content := json.Data
		if bot.Prefix != "" {
			content = bot.Prefix + content
		}
		if bot.Suffix != "" {
			content = content + bot.Suffix
		}
		log.Infof("向 %s发送消息：%s", toUser, content)
		postContent(toUser, content, token)
	}
}

// 定义接收和回复消息的数据结构

type ReplyMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
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
	log.Infof("发送公众号消息结果 %+v", json)
}
