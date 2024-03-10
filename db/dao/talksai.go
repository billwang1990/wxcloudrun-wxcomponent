package dao

import (
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/model"
	"gorm.io/gorm/clause"
)

const talksaiBotTableName = "talks_ai_bot"

// CreateOrUpdateAuthorizerRecord 创建或更新授权账号信息
func CreateOrUpdateTalksAIBot(record *model.TalksAIBot) error {
	var err error
	cli := db.Get()
	if err = cli.Table(talksaiBotTableName).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(record).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
