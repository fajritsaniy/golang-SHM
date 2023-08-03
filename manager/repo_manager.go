package manager

import "github.com/jutionck/golang-db-sinar-harapan-makmur-orm/repository"

type RepositoryManager interface {
	// kumpulan repo disini
	BrandRepo() repository.BrandRepository
	VehicleRepo() repository.VehicleRepository
	CustomerRepo() repository.CustomerRepository
	EmployeeRepo() repository.EmployeeRepository
	TransactionRepo() repository.TransactionRepository
	FileRepo() repository.FileRepository
	UserRepo() repository.UserRepository
}

type repositoryManager struct {
	infra InfraManager
}

func (r *repositoryManager) UserRepo() repository.UserRepository {
	return repository.NewUserRepository(r.infra.Conn())
}

func (r *repositoryManager) FileRepo() repository.FileRepository {
	return repository.NewFileRepository(r.infra.UploadLocation())
}

func (r *repositoryManager) CustomerRepo() repository.CustomerRepository {
	return repository.NewCustomerRepository(r.infra.Conn())
}

func (r *repositoryManager) EmployeeRepo() repository.EmployeeRepository {
	return repository.NewEmployeeRepository(r.infra.Conn())
}

func (r *repositoryManager) TransactionRepo() repository.TransactionRepository {
	return repository.NewTransactionRepository(r.infra.Conn())
}

func (r *repositoryManager) BrandRepo() repository.BrandRepository {
	return repository.NewBrandRepository(r.infra.Conn())
}

func (r *repositoryManager) VehicleRepo() repository.VehicleRepository {
	return repository.NewVehicleRepository(r.infra.Conn())
}

func NewRepositoryManager(infra InfraManager) RepositoryManager {
	return &repositoryManager{infra: infra}
}
