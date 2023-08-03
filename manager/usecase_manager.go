package manager

import (
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/usecase"
)

type UseCaseManager interface {
	BrandUseCase() usecase.BrandUseCase
	VehicleUseCase() usecase.VehicleUseCase
	CustomerUseCase() usecase.CustomerUseCase
	EmployeeUseCase() usecase.EmployeeUseCase
	TransactionUseCase() usecase.TransactionUseCase
	FileUseCase() usecase.FileUseCase
}

type useCaseManager struct {
	repoManager RepositoryManager
}

func (u *useCaseManager) CustomerUseCase() usecase.CustomerUseCase {
	return usecase.NewCustomerUseCase(u.repoManager.CustomerRepo())
}

func (u *useCaseManager) FileUseCase() usecase.FileUseCase {
	return usecase.NewFileUseCase(u.repoManager.FileRepo())
}

func (u *useCaseManager) EmployeeUseCase() usecase.EmployeeUseCase {
	return usecase.NewEmployeeUseCase(u.repoManager.EmployeeRepo())
}

func (u *useCaseManager) TransactionUseCase() usecase.TransactionUseCase {
	return usecase.NewTransactionUseCase(u.repoManager.TransactionRepo(), u.VehicleUseCase(), u.EmployeeUseCase(), u.CustomerUseCase())
}

func (u *useCaseManager) BrandUseCase() usecase.BrandUseCase {
	return usecase.NewBrandUseCase(u.repoManager.BrandRepo())
}

func (u *useCaseManager) VehicleUseCase() usecase.VehicleUseCase {
	return usecase.NewVehicleUseCase(u.repoManager.VehicleRepo(), u.BrandUseCase(), u.FileUseCase())
}

func NewUseCaseManager(repoManager RepositoryManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
