package model

import "time"

type Transaction struct {
	BaseModel
	TransactionDate time.Time `json:"transactionDate"`
	VehicleID       string    `json:"vehicleId"`
	Vehicle         Vehicle   `gorm:"foreignKey:VehicleID" json:"vehicle"`
	CustomerID      string    `json:"customerId"`
	Customer        Customer  `gorm:"foreignKey:CustomerID" json:"customer"`
	EmployeeID      string    `json:"employeeId"`
	Employee        Employee  `gorm:"foreignKey:EmployeeID" json:"employee"`
	Type            string    `gorm:"check:type IN ('online', 'offline')" json:"type"`
	Qty             int       `json:"qty"`
	PaymentAmount   int64     `json:"paymentAmount"`
}

func (t *Transaction) IsValidType() bool {
	return t.Type == "online" || t.Type == "offline"
}

func (Transaction) TableName() string {
	return "trx_transaction"
}
