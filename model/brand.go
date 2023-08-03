package model

import validation "github.com/go-ozzo/ozzo-validation"

type Brand struct {
	BaseModel
	Name     string    `json:"name"`
	Vehicles []Vehicle `json:"vehicles,omitempty"`
}

func (b Brand) TableName() string {
	return "mst_brand"
}

func (b Brand) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Required, validation.Length(3, 30)),
	)
}
