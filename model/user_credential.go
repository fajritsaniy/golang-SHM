package model

type UserCredential struct {
	BaseModel
	UserName string `gorm:"unique;size:50;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	IsActive bool   `gorm:"default:true"`
}

func (UserCredential) TableName() string {
	return "mst_user"
}
