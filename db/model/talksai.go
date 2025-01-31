package model

type TalksAIBot struct {
	BotID          string `gorm:"column:botid" json:"botid"`
	Verified       bool   `gorm:"column:verified" json:"verified"`
	AppID          string `gorm:"column:appid" json:"appid"`
	Filters        string `gorm:"column:filters" json:"filters"`
	ExcludeFilters string `gorm:"column:excludefilters" json:"excludefilters"`
	Prefix         string `gorm:"column:prefix" json:"prefix"`
	Suffix         string `gorm:"column:suffix" json:"suffix"`
	//公众号名称
	Name string `gorm:"column:name" json:"name"`
}
