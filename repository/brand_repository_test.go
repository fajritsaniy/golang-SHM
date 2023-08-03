package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

var brandDummies = []model.Brand{
	{
		BaseModel: model.BaseModel{ID: "1", CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
		Name:      "Honda",
	},
	{
		BaseModel: model.BaseModel{ID: "2", CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
		Name:      "Toyota",
	},
	{
		BaseModel: model.BaseModel{ID: "3", CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
		Name:      "BMW",
	},
}

type BrandRepoTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *BrandRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.mock = mock
	dialect := postgres.New(postgres.Config{
		Conn: db,
	})
	suite.DB, err = gorm.Open(dialect)
	assert.NoError(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestGetAllBrand_Success() {
	brandRowDummies := make([]model.Brand, len(brandDummies))
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
	for i, brand := range brandDummies {
		brandRowDummies[i] = brand
		rows.AddRow(brand.ID, brand.Name, brand.CreatedAt, brand.UpdatedAt)
	}
	expectedQuery := `SELECT \* FROM "mst_brand"`
	suite.mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
	repo := NewBrandRepository(suite.DB)
	listBrand, err := repo.List()
	assert.Equal(suite.T(), brandRowDummies, listBrand)
	assert.NoError(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestGetAllMenu_DBErrorFail() {
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
	for _, brand := range brandDummies {
		rows.AddRow(brand.ID, brand.Name, brand.CreatedAt, brand.UpdatedAt)
	}
	expectedQuery := `SELECT \* FROM "mst_brand"`
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("db error"))
	repo := NewBrandRepository(suite.DB)
	listMenu, err := repo.List()
	assert.Nil(suite.T(), listMenu)
	assert.Error(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestGetById_Success() {
	brandDm := &brandDummies[0]
	brandRow := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(brandDm.ID, brandDm.Name, brandDm.CreatedAt, brandDm.UpdatedAt)
	expectedQuery := `SELECT \* FROM "mst_brand" WHERE id(\s*)=(\s*)\$1 AND "mst_brand"."deleted_at" IS NULL ORDER BY "mst_brand"."id" LIMIT 1`
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(brandDm.ID).WillReturnRows(brandRow)
	repo := NewBrandRepository(suite.DB)
	brand, err := repo.Get(brandDm.ID)
	assert.Equal(suite.T(), *brandDm, *brand)
	assert.NoError(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestGetById_DBErrorFail() {
	brandDm := brandDummies[0]
	brandRow := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
	brandRow.AddRow(brandDm.ID, brandDm.Name, brandDm.CreatedAt, brandDm.UpdatedAt)
	expectedQuery := `SELECT \* FROM "mst_brand" WHERE id(\s*)=(\s*)\$1 AND "mst_brand"."deleted_at" IS NULL ORDER BY "mst_brand"."id" LIMIT 1`
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(brandDm.ID).WillReturnError(errors.New("db error"))
	repo := NewBrandRepository(suite.DB)
	brand, err := repo.Get(brandDm.ID)
	assert.Nil(suite.T(), brand)
	assert.Error(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestSearchBrand_Success() {
	brandRowDummies := make([]model.Brand, len(brandDummies))
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
	for i, brand := range brandDummies {
		brandRowDummies[i] = brand
		rows.AddRow(brand.ID, brand.Name, brand.CreatedAt, brand.UpdatedAt)
	}
	expectedQuery := `SELECT \* FROM "mst_brand" WHERE \"name\"(\s*)=(\s*)\$1 AND "mst_brand"."deleted_at" IS NULL`
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs("Honda").WillReturnRows(rows)
	repo := NewBrandRepository(suite.DB)
	filter := map[string]interface{}{"name": "Honda"}
	listBrand, err := repo.Search(filter)
	assert.Equal(suite.T(), brandRowDummies, listBrand)
	assert.NoError(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestSearchBrand_DBErrorFail() {
	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
	for _, brand := range brandDummies {
		rows.AddRow(brand.ID, brand.Name, brand.CreatedAt, brand.UpdatedAt)
	}
	expectedQuery := `SELECT \* FROM "mst_brand" WHERE \"name\"(\s*)=(\s*)\$1 AND "mst_brand"."deleted_at" IS NULL`
	suite.mock.ExpectQuery(expectedQuery).WillReturnError(errors.New("db error"))
	repo := NewBrandRepository(suite.DB)
	filter := map[string]interface{}{"name": "Honda"}
	listMenu, err := repo.Search(filter)
	assert.Nil(suite.T(), listMenu)
	assert.Error(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestDeleteMenu_Success() {
	suite.mock.ExpectBegin()
	expectedQuery := `UPDATE "mst_brand"`
	suite.mock.ExpectExec(expectedQuery).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()
	repo := NewBrandRepository(suite.DB)
	err := repo.Delete("1")
	assert.Nil(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestDeleteMenu_DBErrorFail() {
	expectedQuery := `UPDATE "mst_brand"`
	suite.mock.ExpectExec(expectedQuery).
		WillReturnError(errors.New("db error"))
	repo := NewBrandRepository(suite.DB)
	err := repo.Delete("1")
	assert.Error(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestCountByName_Success() {
	brandDm := brandDummies[0]
	filter := map[string]interface{}{"name": "Honda"}
	brandRow := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(brandDm.ID, brandDm.Name, brandDm.CreatedAt, brandDm.UpdatedAt)
	expectedQuery := `SELECT \* FROM "mst_brand" WHERE name(\s*)=(\s*)\$1 AND "mst_brand"."deleted_at" IS NULL ORDER BY "mst_brand"."id" LIMIT 1`
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(filter).WillReturnRows(brandRow)
	repo := NewBrandRepository(suite.DB)
	brand, err := repo.Search(filter)
	assert.Equal(suite.T(), brandDm, brand)
	assert.NoError(suite.T(), err)
}

func (suite *BrandRepoTestSuite) TestCountByName_DBErrorFail() {
	brandDm := brandDummies[0]
	filter := map[string]interface{}{"name": "Honda"}
	brandRow := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
	brandRow.AddRow(brandDm.ID, brandDm.Name, brandDm.CreatedAt, brandDm.UpdatedAt)
	expectedQuery := `SELECT \* FROM "mst_brand" WHERE name(\s*)=(\s*)\$1 AND "mst_brand"."deleted_at" IS NULL`
	suite.mock.ExpectQuery(expectedQuery).
		WithArgs(brandDm.ID).WillReturnError(errors.New("db error"))
	repo := NewBrandRepository(suite.DB)
	brand, err := repo.Search(filter)
	assert.Nil(suite.T(), brand)
	assert.Error(suite.T(), err)
}

func TestBrandRepoTestSuite(t *testing.T) {
	suite.Run(t, new(BrandRepoTestSuite))
}
