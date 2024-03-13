package model

type CacheNewAuthRecord struct {
	AuthorizerAppid   string `gorm:"column:appid" json:"appid"`
	AuthorizerAppName string `gorm:"column:appname" json:"appname"`
	AuthorizationCode string `gorm:"column:appcode" json:"appcode"`
}
