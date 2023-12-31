package usecase

import (
	"fmt"

	"github.com/fajritsaniy/golang-SHM/model"
	"github.com/fajritsaniy/golang-SHM/repository"
	"github.com/fajritsaniy/golang-SHM/utils"
)

type CustomerUseCase interface {
	BaseUseCase[model.Customer]
	BaseUseCaseEmailPhone[model.Customer]
	AppendCustomerVehicle(payload *model.Customer, association any) error
}

type customerUseCase struct {
	repo repository.CustomerRepository
}

func CustomerNotFoundMessage(id string) string {
	return fmt.Sprintf("customers with ID %s not found", id)
}

func (c *customerUseCase) DeleteData(id string) error {
	customer, err := c.FindById(id)
	if err != nil {
		return fmt.Errorf(CustomerNotFoundMessage(id))
	}
	return c.repo.Delete(customer.ID)
}

func (c *customerUseCase) FindAll() ([]model.Customer, error) {
	return c.repo.List()
}

func (c *customerUseCase) FindById(id string) (*model.Customer, error) {
	customer, err := c.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf(CustomerNotFoundMessage(id))
	}
	return customer, nil
}

func (c *customerUseCase) SaveData(payload *model.Customer) error {
	if payload.ID != "" {
		_, err := c.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf(CustomerNotFoundMessage(payload.ID))
		}
	}

	// create user credential (recommended use transactional)
	password, err := utils.HashPassword("password")
	if err != nil {
		return err
	}
	userCredential := model.UserCredential{
		UserName: payload.Email,
		Password: password,
		IsActive: false,
	}
	payload.UserCredential = userCredential
	return c.repo.Save(payload)
}

func (c *customerUseCase) SearchBy(by map[string]interface{}) ([]model.Customer, error) {
	customers, err := c.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}
	return customers, nil
}

func (c *customerUseCase) FindByEmail(email string) (*model.Customer, error) {
	customer, err := c.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("customer with email %s not found", email)
	}
	return customer, nil
}

func (c *customerUseCase) FindByPhone(phone string) (*model.Customer, error) {
	customer, err := c.repo.GetByPhone(phone)
	if err != nil {
		return nil, fmt.Errorf("customer with phone number %s not found", phone)
	}
	return customer, nil
}

func (c *customerUseCase) AppendCustomerVehicle(payload *model.Customer, association any) error {
	return c.repo.CreateCustomerVehicle(payload, association)
}

func NewCustomerUseCase(repo repository.CustomerRepository) CustomerUseCase {
	return &customerUseCase{repo: repo}
}
