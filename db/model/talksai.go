package model


type TalksAIBot struct {
	BotID    string   `gorm:"column:botid" json:"botid"`
	AppID    string   `gorm:"column:botid" json:"appid"`
	Filters  []string `gorm:"column:filters" json:"filters"`
	Prefix   string   `gorm:"column:prefix" json:"prefix"`
	Suffix   string   `gorm:"column:suffix" json:"suffix"`
}
