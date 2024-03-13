package model

type TalksAIBot struct {
	BotID   string `gorm:"column:botid" json:"botid"`
	AppID   string `gorm:"column:appid" json:"appid"`
	Filters string `gorm:"column:filters" json:"filters"`
	Prefix  string `gorm:"column:prefix" json:"prefix"`
	Suffix  string `gorm:"column:suffix" json:"suffix"`
	//公众号名称
	Name  string `gorm:"column:name" json:"name"` 
}
