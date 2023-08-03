package model

import "time"

type Employee struct {
	BaseModel
	FirstName        string    `gorm:"size:30" json:"firstName"`
	LastName         string    `gorm:"size:30" json:"lastName"`
	Address          string    `json:"address"`
	Email            string    `gorm:"unique;size:30" json:"email"`
	PhoneNumber      string    `gorm:"unique;size:15" json:"phoneNumber"`
	Bod              time.Time `json:"bod"`
	Position         string    `json:"position"`
	Salary           int64     `gorm:"default:0" json:"salary"`
	ManagerID        *string   `json:"managerID"`
	Manager          *Employee `gorm:"foreignKey:ManagerID"`
	UserCredentialID string
	UserCredential   UserCredential `gorm:"foreignKey:UserCredentialID;unique"`
}

func (Employee) TableName() string {
	return "mst_employee"
}
