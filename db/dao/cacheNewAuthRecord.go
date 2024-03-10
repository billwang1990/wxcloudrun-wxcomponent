package dao

import (
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db"
	"gorm.io/gorm/clause"
)

const cacheTmpAppidWithCode = "cacheTmpAppidWithCode"

// WxCallbackRule 回调消息转发规则
type CacheNewAuthRecord struct {
	AuthorizerAppid   string `gorm:"column:appid" json:"appid"`
	AuthorizationCode string `gorm:"column:appcode" json:"appcode"`
}

// CreateOrUpdateAuthorizerRecord 创建或更新授权账号信息
func CreateOrUpdateAuthorizerAppWithCode(appid string, authCode string) error {
	record := &CacheNewAuthRecord{AuthorizerAppid: appid, AuthorizationCode: authCode}
	var err error
	cli := db.Get()
	if err = cli.Table(cacheTmpAppidWithCode).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(record).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}