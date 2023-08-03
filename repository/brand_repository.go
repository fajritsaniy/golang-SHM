package repository

import (
	"fmt"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/utils/common"
	"gorm.io/gorm"
)

type BrandRepository interface {
	BaseRepository[model.Brand]
	BaseRepositoryPaging[model.Brand]
	CountByName(name string, id string) (int64, error)
}

type brandRepository struct {
	db *gorm.DB
}

func (b *brandRepository) Delete(id string) error {
	return b.db.Delete(&model.Brand{}, "id=?", id).Error
}

func (b *brandRepository) Get(id string) (*model.Brand, error) {
	var brand model.Brand
	result := b.db.First(&brand, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &brand, nil
}

func (b *brandRepository) List() ([]model.Brand, error) {
	var brands []model.Brand
	result := b.db.Find(&brands).Error
	if result != nil {
		return nil, result
	}
	return brands, nil
}

func (b *brandRepository) Save(payload *model.Brand) error {
	return b.db.Save(payload).Error
}

func (b *brandRepository) Search(by map[string]interface{}) ([]model.Brand, error) {
	var brands []model.Brand
	result := b.db.Where(by).Find(&brands).Error
	if result != nil {
		return nil, result
	}
	return brands, nil
}

func (b *brandRepository) CountByName(name string, id string) (int64, error) {
	var count int64
	var result *gorm.DB
	if id != "" {
		result = b.db.Model(&model.Brand{}).Where("name ILIKE ? AND id <> ?", "%"+name+"%", id).Count(&count)
	} else {
		result = b.db.Model(&model.Brand{}).Where("name ILIKE ?", "%"+name+"%").Count(&count)
	}
	if err := result.Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (b *brandRepository) Paging(requestQueryParams dto.RequestQueryParams) ([]model.Brand, dto.Paging, error) {
	paginationQuery, orderQuery := b.pagingValidate(requestQueryParams)
	var brands []model.Brand
	result := b.db.Preload("Vehicles").Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&brands).Error
	if result != nil {
		return nil, dto.Paging{}, result
	}
	var totalRows int64
	result = b.db.Model(&model.Brand{}).Count(&totalRows).Error
	if result != nil {
		return nil, dto.Paging{}, result
	}
	return brands, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func (b *brandRepository) pagingValidate(requestQueryParams dto.RequestQueryParams) (dto.PaginationQuery, string) {
	var paginationQuery dto.PaginationQuery
	paginationQuery = common.GetPaginationParams(requestQueryParams.PaginationParam)
	orderQuery := "id"
	if requestQueryParams.QueryParams.Order != "" && requestQueryParams.QueryParams.Sort != "" {
		sorting := "ASC"
		if requestQueryParams.QueryParams.Sort == "desc" {
			sorting = "DESC"
		}
		orderQuery = fmt.Sprintf("%s %s", requestQueryParams.QueryParams.Order, sorting)
	}
	return paginationQuery, orderQuery
}

func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{db: db}
}
