package models

type Client struct {
	BaseModel
	ClientName   string `json:"clientName" binding:"required" gorm:"type:varchar(100);not null"`
	ClientKey    string `json:"clientKey" gorm:"type:varchar(255);uniqueIndex;not null"`
	ClientSecret string `json:"clientSecret" gorm:"type:varchar(255);not null"`
	IsActive     bool   `json:"isActive" gorm:"default:true"`
}
