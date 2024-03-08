package wxcallback

import (
	"io/ioutil"
	"net/http"
	"time"
	"bytes"
    "fmt"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"

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

func postContent(content  string, token string) {
	log.Infof("wyq-------// postContent %s", content)
	    // 定义要发送的JSON数据
		jsonData := []byte(`{
			"touser": "opnbu552g7sy8s63dgm-M60lg7Og",
			"msgtype": "text",
			"text": {
				"content": "\n—————保罗AI客服回复"
			}
		}`)
	
		// 创建POST请求
		url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + token
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("发送消息 失败 Error:", err)
			return
		}
	
		// 设置请求头
		req.Header.Set("Content-Type", "application/json")
	
		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
	
		// 读取响应数据
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	
		// 打印响应数据
		fmt.Println(string(body))
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
	if err := dao.AddBizCallBackRecord(&r); err != nil {
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}
	mytoken := "78_uCzSfzKazw5oy3atpp546ZuEWQg0GbrCptn6E-Gv09QKR_oNji0-o6ey7u-RxLsR4Loq17LAXxF1d-83uYCZ4IG2ZjNykaBGzRSu3gC6aNGox_JTy9NHo4KgvACvJvNKI3w-5bC3-ihC72p2WPZdAGDGZJ"
	postContent("", mytoken)
	log.Infof("wyq-------// 转发到用户配置的地址")
	// 转发到用户配置的地址
	proxyOpen, err := proxyCallbackMsg("", json.MsgType, json.Event, string(body), c)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}
	if !proxyOpen {
		log.Infof("wyq-------// proxyOpen")
		c.String(http.StatusOK, "success")
	}
}
