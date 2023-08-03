package repository

import (
	"fmt"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	BaseRepository[model.Customer]
	ListCustomerUser() ([]model.Customer, error)
	GetByUser(userId string) (*model.Customer, error)
	BaseRepositoryEmailPhone[model.Customer]
	CreateCustomerVehicle(payload *model.Customer, association any) error
}

type customerRepository struct {
	db *gorm.DB
}

func (c *customerRepository) Search(by map[string]interface{}) ([]model.Customer, error) {
	var customers []model.Customer
	result := c.db.Where(by).Find(&customers).Error
	if result != nil {
		return nil, result
	}
	return customers, nil
}

func (c *customerRepository) List() ([]model.Customer, error) {
	var customers []model.Customer
	result := c.db.Find(&customers).Error
	if result != nil {
		return nil, result
	}
	return customers, nil
}

func (c *customerRepository) Get(id string) (*model.Customer, error) {
	var customer model.Customer
	result := c.db.First(&customer, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &customer, nil
}

func (c *customerRepository) ListCustomerUser() ([]model.Customer, error) {
	var customers []model.Customer
	result := c.db.Preload("UserCredential").Order("created_at").Find(&customers).Error
	if result != nil {
		return nil, result
	}

	return customers, nil
}

func (c *customerRepository) GetByUser(userId string) (*model.Customer, error) {
	var customer model.Customer
	result := c.db.Preload("UserCredential").First(&customer, "user_credential_id=?", userId).Error
	if result != nil {
		return nil, result
	}

	return &customer, nil
}

func (c *customerRepository) Save(payload *model.Customer) error {
	return c.db.Save(payload).Error
}

func (c *customerRepository) Delete(id string) error {
	return c.db.Delete(&model.Customer{}, "id=?", id).Error
}

func (c *customerRepository) GetByEmail(email string) (*model.Customer, error) {
	var customer model.Customer
	result := c.db.Select("id, email").First(&customer, "email=?", email).Error
	if result != nil {
		return nil, result
	}
	return &customer, nil
}

func (c *customerRepository) GetByPhone(phone string) (*model.Customer, error) {
	var customer model.Customer
	result := c.db.Select("id, phone_number").First(&customer, "phone_number=?", phone).Error
	if result != nil {
		return nil, result
	}
	return &customer, nil
}

func (c *customerRepository) CreateCustomerVehicle(payload *model.Customer, association interface{}) error {
	vehicle, ok := association.(*model.Vehicle)
	if !ok {
		return fmt.Errorf("invalid vehicle association")
	}

	if err := c.db.Model(vehicle).Association("Customers").Append(payload); err != nil {
		return err
	}

	return nil
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}
