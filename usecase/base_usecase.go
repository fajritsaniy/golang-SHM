package usecase

import "github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"

type BaseUseCase[T any] interface {
	SearchBy(by map[string]interface{}) ([]T, error)
	FindAll() ([]T, error)
	FindById(id string) (*T, error)
	SaveData(payload *T) error
	DeleteData(id string) error
}

type BaseUseCaseEmailPhone[T any] interface {
	FindByEmail(email string) (*T, error)
	FindByPhone(phone string) (*T, error)
}

type BaseUseCasePaging[T any] interface {
	Pagination(requestQueryParams dto.RequestQueryParams) ([]T, dto.Paging, error)
}
