package model

import "time"

type Customer struct {
	BaseModel
	FirstName        string    `gorm:"size:30" json:"firstName"`
	LastName         string    `gorm:"size:30" json:"lastName"`
	Address          string    `json:"address"`
	Email            string    `gorm:"unique;size:30" json:"email"`
	PhoneNumber      string    `gorm:"unique;size:15" json:"phoneNumber"`
	Bod              time.Time `json:"bod"`
	UserCredentialID string
	UserCredential   UserCredential `gorm:"foreignKey:UserCredentialID"`
	Vehicles         []Vehicle      `gorm:"many2many:customer_vehicles;"`
}

func (Customer) TableName() string {
	return "mst_customer"
}
