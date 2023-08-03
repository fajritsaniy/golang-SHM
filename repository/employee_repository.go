package repository

import (
	"errors"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	BaseRepository[model.Employee]
	ListEmployeeUser() ([]model.Employee, error)
	GetByUser(userId string) (*model.Employee, error)
	ListEmployeeByManager(managerId string) ([]model.Employee, error)
	BaseRepositoryEmailPhone[model.Employee]
}

type employeeRepository struct {
	db *gorm.DB
}

func (e *employeeRepository) Search(by map[string]interface{}) ([]model.Employee, error) {
	var employees []model.Employee
	result := e.db.Where(by).Find(&employees).Error
	if result != nil {
		return nil, result
	}
	return employees, nil
}

func (e *employeeRepository) List() ([]model.Employee, error) {
	var employees []model.Employee
	result := e.db.Preload("Manager").Find(&employees).Error
	if result != nil {
		return nil, result
	}
	return employees, nil
}

func (e *employeeRepository) Get(id string) (*model.Employee, error) {
	var employee model.Employee
	result := e.db.First(&employee, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &employee, nil
}

func (e *employeeRepository) ListEmployeeUser() ([]model.Employee, error) {
	var employees []model.Employee
	result := e.db.Preload("UserCredential").Order("created_at").Find(&employees).Error
	if result != nil {
		return nil, result
	}

	return employees, nil
}

func (e *employeeRepository) GetByUser(userId string) (*model.Employee, error) {
	var employee model.Employee
	result := e.db.Preload("UserCredential").Where("user_credential_id = ?", userId).First(&employee).Error
	if result != nil {
		return nil, result
	}

	return &employee, nil
}

func (e *employeeRepository) Save(payload *model.Employee) error {
	return e.db.Save(payload).Error
}

func (e *employeeRepository) Delete(id string) error {
	return e.db.Delete(&model.Employee{}, "id=?", id).Error
}

func (e *employeeRepository) GetByEmail(email string) (*model.Employee, error) {
	var employee model.Employee
	err := e.db.Where("email = ?", email).Select("id, email").First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &employee, nil
}

func (e *employeeRepository) GetByPhone(phone string) (*model.Employee, error) {
	var employee model.Employee
	err := e.db.Where("phone_number = ?", phone).Select("id, phone_number").First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &employee, nil
}

func (e *employeeRepository) ListEmployeeByManager(managerId string) ([]model.Employee, error) {
	var employees []model.Employee
	result := e.db.Preload("Manager").Where("manager_id = ?", managerId).Find(&employees).Error
	if result != nil {
		return nil, result
	}
	return employees, nil
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}
