package usecase

import (
	"fmt"
	"time"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/repository"
)

type TransactionUseCase interface {
	RegisterNewTransaction(payload *model.Transaction) error
	FindAllTransaction() ([]model.Transaction, error)
	FindByTransaction(id string) (model.Transaction, error)
}

type transactionUseCase struct {
	repo       repository.TransactionRepository
	vehicleUC  VehicleUseCase
	employeeUC EmployeeUseCase
	customerUC CustomerUseCase
}

func (t *transactionUseCase) RegisterNewTransaction(payload *model.Transaction) error {
	// Recommended use Transaction
	// get vehicle
	vehicle, err := t.vehicleUC.FindById(payload.VehicleID)
	if err != nil {
		return err
	}

	// get employee
	employee, err := t.employeeUC.FindById(payload.EmployeeID)
	if err != nil {
		return err
	}

	// get customer
	customer, err := t.customerUC.FindById(payload.CustomerID)
	if err != nil {
		return err
	}

	// append customer vehicle
	err = t.customerUC.AppendCustomerVehicle(customer, vehicle)
	if err != nil {
		return fmt.Errorf("failed to append customer vehicle")
	}

	// validate stock
	if vehicle.Stock < payload.Qty {
		return fmt.Errorf("not enough stock")
	}

	//update stock
	err = t.vehicleUC.UpdateVehicleStock(payload.Qty, vehicle.ID)
	if err != nil {
		return fmt.Errorf("failed to update stock")
	}

	payload.Vehicle = *vehicle
	payload.Customer = *customer
	payload.Employee = *employee
	payload.TransactionDate = time.Now()
	payload.PaymentAmount = int64(vehicle.SalePrice)

	err = t.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	return nil
}

func (t *transactionUseCase) FindAllTransaction() ([]model.Transaction, error) {
	return t.repo.List()
}

func (t *transactionUseCase) FindByTransaction(id string) (model.Transaction, error) {
	return t.repo.Get(id)
}

func NewTransactionUseCase(
	repo repository.TransactionRepository,
	vehicleUC VehicleUseCase,
	employeeUC EmployeeUseCase,
	customerUC CustomerUseCase) TransactionUseCase {
	return &transactionUseCase{
		repo:       repo,
		vehicleUC:  vehicleUC,
		employeeUC: employeeUC,
		customerUC: customerUC,
	}
}
