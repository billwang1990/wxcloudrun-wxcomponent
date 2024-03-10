package model

type CacheNewAuthRecord struct {
	AuthorizerAppid   string `gorm:"column:appid" json:"appid"`
	AuthorizationCode string `gorm:"column:appcode" json:"appcode"`
}
