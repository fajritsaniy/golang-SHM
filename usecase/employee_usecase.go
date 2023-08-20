package usecase

import (
	"fmt"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/utils"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/repository"
)

type EmployeeUseCase interface {
	BaseUseCase[model.Employee]
	BaseUseCaseEmailPhone[model.Employee]
	FindAllEmployeeByManager(managerId string) ([]model.Employee, error)
}

type employeeUseCase struct {
	repo repository.EmployeeRepository
}

func employeeIDNotFoundMessage(id string) string {
	return fmt.Sprintf("employee with ID %s not found", id)
}

func (e *employeeUseCase) DeleteData(id string) error {
	employee, err := e.FindById(id)
	if err != nil {
		return fmt.Errorf(employeeIDNotFoundMessage(id))
	}
	return e.repo.Delete(employee.ID)
}

func (e *employeeUseCase) FindAll() ([]model.Employee, error) {
	return e.repo.List()
}

func (e *employeeUseCase) FindById(id string) (*model.Employee, error) {
	employee, err := e.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf(employeeIDNotFoundMessage(id))
	}
	return employee, nil
}

func (e *employeeUseCase) SaveData(payload *model.Employee) error {
	if payload.ID != "" {
		_, err := e.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf(employeeIDNotFoundMessage(payload.ID))
		}
	}

	isEmailExist, _ := e.FindByEmail(payload.Email)
	if isEmailExist != nil && isEmailExist.Email == payload.Email {
		return fmt.Errorf("employee with email: %v exists", payload.Email)
	}

	isPhoneNumberExist, _ := e.FindByPhone(payload.PhoneNumber)
	if isPhoneNumberExist != nil && isPhoneNumberExist.PhoneNumber == payload.PhoneNumber {
		return fmt.Errorf("employee with phone number: %v exists", payload.PhoneNumber)
	}

	if payload.ManagerID != nil {
		manager, _ := e.FindById(*payload.ManagerID)
		payload.Manager = manager
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
	return e.repo.Save(payload)
}

func (e *employeeUseCase) SearchBy(by map[string]interface{}) ([]model.Employee, error) {
	employees, err := e.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}
	return employees, nil
}

func (e *employeeUseCase) FindByEmail(email string) (*model.Employee, error) {
	employee, err := e.repo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("employee with email %s not found", email)
	}
	return employee, nil
}

func (e *employeeUseCase) FindByPhone(phone string) (*model.Employee, error) {
	employee, err := e.repo.GetByPhone(phone)
	if err != nil {
		return nil, fmt.Errorf("employee with phone number %s not found", phone)
	}
	return employee, nil
}

func (e *employeeUseCase) FindAllEmployeeByManager(managerId string) ([]model.Employee, error) {
	return e.repo.ListEmployeeByManager(managerId)
}

func NewEmployeeUseCase(repo repository.EmployeeRepository) EmployeeUseCase {
	return &employeeUseCase{repo: repo}
}
