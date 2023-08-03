package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	BrandID        string     `json:"brandId"`
	Brand          Brand      `json:"brand"`
	Model          string     `gorm:"varchar;size:30" json:"model"`
	ProductionYear int        `gorm:"size:4" json:"productionYear"`
	Color          string     `gorm:"varchar;size:30" json:"color"`
	IsAutomatic    bool       `json:"isAutomatic"`
	Stock          int        `gorm:"check:stock >= 0" json:"stock"`
	SalePrice      int        `gorm:"check:sale_price > 0" json:"salePrice"`
	Status         string     `gorm:"check:status IN ('baru', 'bekas')" json:"status"`
	Customers      []Customer `gorm:"many2many:customer_vehicles;" json:"customers,omitempty"`
	ImgPath        string     `json:"imgPath,omitempty"`
	UrlPath        string     `json:"urlPath"`
	BaseModel
}

func (v *Vehicle) TableName() string {
	return "mst_vehicle"
}

func (v *Vehicle) IsValidStatus() bool {
	return v.Status == "baru" || v.Status == "bekas"
}

func (v *Vehicle) BeforeCreate(tx *gorm.DB) error {
	v.ID = uuid.New().String()
	return nil
}
