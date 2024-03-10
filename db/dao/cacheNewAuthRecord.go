package dao

import (
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/model"
	"gorm.io/gorm/clause"
)

const cacheTmpAppidWithCode = "cache_appid_code"

// CreateOrUpdateAuthorizerRecord 创建或更新授权账号信息
func CreateOrUpdateAuthorizerAppWithCode(appid string, authCode string) error {
	record := &model.CacheNewAuthRecord{AuthorizerAppid: appid, AuthorizationCode: authCode}
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

func GetAuthRecordByCode(code string) (*model.CacheNewAuthRecord, error) {
	var err error
	var kv model.CacheNewAuthRecord
	cli := db.Get()
	if err = cli.Table(cacheTmpAppidWithCode).Where("`appcode` = ?", code).Take(&kv).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &kv, nil
}