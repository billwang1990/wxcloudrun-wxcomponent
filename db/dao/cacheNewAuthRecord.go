package dao

import (
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db"
	"gorm.io/gorm/clause"
)

const cacheTmpAppidWithCode = "cacheTmpAppidWithCode"

type newAuthRecord struct {
	CreateTime                   int64  `json:"CreateTime"`
	AuthorizerAppid              string `json:"AuthorizerAppid"`
	AuthorizationCode            string `json:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int64  `json:"AuthorizationCodeExpiredTime"`
}

// CreateOrUpdateAuthorizerRecord 创建或更新授权账号信息
func CreateOrUpdateAuthorizerAppWithCode(record *newAuthRecord) error {
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

