package repository

import (
	"errors"
	"fmt"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/utils/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VehicleRepository interface {
	BaseRepository[model.Vehicle]
	BaseRepositoryPaging[model.Vehicle]
	UpdateStock(count int, id string) error
}

type vehicleRepository struct {
	db *gorm.DB
}

func (v *vehicleRepository) Search(by map[string]interface{}) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	result := v.db.Where(by).Find(&vehicles)
	if err := result.Error; err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

func (v *vehicleRepository) List() ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	result := v.db.Preload(clause.Associations).Find(&vehicles)
	if err := result.Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (v *vehicleRepository) Get(id string) (*model.Vehicle, error) {
	var vehicle model.Vehicle
	result := v.db.Preload(clause.Associations).First(&vehicle, "id = ?", id)
	if err := result.Error; err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (v *vehicleRepository) Save(payload *model.Vehicle) error {
	result := v.db.Save(payload)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (v *vehicleRepository) Delete(id string) error {
	return v.db.Delete(&model.Vehicle{}, "id=?", id).Error
}

func (v *vehicleRepository) UpdateStock(count int, id string) error {
	result := v.db.Model(&model.Vehicle{}).Where("id=?", id).Update("stock", gorm.Expr("stock - ?", count))
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (v *vehicleRepository) Paging(requestQueryParams dto.RequestQueryParams) ([]model.Vehicle, dto.Paging, error) {
	var paginationQuery dto.PaginationQuery
	var vehicles []model.Vehicle
	paginationQuery = common.GetPaginationParams(requestQueryParams.PaginationParam)
	orderQuery := "id"
	if requestQueryParams.QueryParams.Order != "" && requestQueryParams.QueryParams.Sort != "" {
		sorting := "ASC"
		if requestQueryParams.QueryParams.Sort == "desc" {
			sorting = "DESC"
		}
		orderQuery = fmt.Sprintf("%s %s", requestQueryParams.QueryParams.Order, sorting)
	}

	res := v.db.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Preload(clause.Associations).Find(&vehicles)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.Paging{}, nil
		} else {
			return nil, dto.Paging{}, err
		}
	}

	var totalRows int64
	err := v.db.Model(&model.Vehicle{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return vehicles, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewVehicleRepository(db *gorm.DB) VehicleRepository {
	return &vehicleRepository{db: db}
}
