package usecase

import (
	"errors"
	"testing"

	"github.com/fajritsaniy/golang-SHM/model"
	"github.com/fajritsaniy/golang-SHM/model/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// brandDummies data for mock
var brandDummies = []model.Brand{
	{
		BaseModel: model.BaseModel{ID: "1"},
		Name:      "Honda",
	},
	{
		BaseModel: model.BaseModel{ID: "2"},
		Name:      "Toyota",
	},
	{
		BaseModel: model.BaseModel{ID: "3"},
		Name:      "BMW",
	},
}

const repositoryErrorMessage = "Error on Repository!"

// repoMock as repository mock, because use case need repo for running
type repoMock struct {
	mock.Mock
}

// Setup all repository here (mock)
func (r *repoMock) Delete(id string) error {
	ret := r.Called(id)
	return ret.Error(0)
}

func (r *repoMock) Get(id string) (*model.Brand, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Brand), nil
}

func (r *repoMock) List() ([]model.Brand, error) {
	args := r.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Brand), nil
}

func (r *repoMock) Save(payload *model.Brand) error {
	ret := r.Called(payload)
	return ret.Error(0)
}

func (r *repoMock) Search(by map[string]interface{}) ([]model.Brand, error) {
	args := r.Called(by)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Brand), nil
}

func (r *repoMock) CountByName(name string, id string) (int64, error) {
	args := r.Called(name, id)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), nil
}

func (b *repoMock) Paging(requestQueryParams dto.RequestQueryParams) ([]model.Brand, dto.Paging, error) {
	args := b.Called(requestQueryParams)
	return args.Get(0).([]model.Brand), args.Get(1).(dto.Paging), args.Error(2)
}

func (suite *BrandUseCaseTestSuite) TestIsNameExistsSuccess() {
	var countBrand int64 = 0
	suite.repoMock.On("CountByName", "Honda", "1").Return(countBrand, nil)
	useCase := NewBrandUseCase(suite.repoMock)
	count, err := useCase.IsNameExists("Honda", "1")
	assert.Equal(suite.T(), false, count)
	assert.Nil(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestIsNameExistsRepoErrorFail() {
	var countBrand int64 = 1
	suite.repoMock.On("CountByName", "Honda", "1").Return(countBrand, errors.New(repositoryErrorMessage))
	useCase := NewBrandUseCase(suite.repoMock)
	count, err := useCase.IsNameExists("Honda", "1")
	assert.Equal(suite.T(), true, count)
	assert.Error(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestFindAllSuccess() {
	suite.repoMock.On("List").Return(brandDummies, nil)
	useCase := NewBrandUseCase(suite.repoMock)
	brands, err := useCase.FindAll()
	assert.Equal(suite.T(), brandDummies, brands)
	assert.Nil(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestFindAllRepoErrorFail() {
	suite.repoMock.On("List").Return(nil, errors.New(repositoryErrorMessage))
	useCase := NewBrandUseCase(suite.repoMock)
	list, err := useCase.FindAll()
	assert.Nil(suite.T(), list)
	assert.Error(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestDeleteBrandSuccess() {
	suite.repoMock.On("Get", "1").Return(&brandDummies[0], nil)
	suite.repoMock.On("Delete", "1").Return(nil)
	useCase := NewBrandUseCase(suite.repoMock)
	err := useCase.DeleteData("1")
	assert.Nil(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestDeleteBrandRepoErrorFail() {
	suite.repoMock.On("Get", "1").Return(nil, errors.New(repositoryErrorMessage))
	suite.repoMock.On("Delete", "1").Return(errors.New(repositoryErrorMessage))
	useCase := NewBrandUseCase(suite.repoMock)
	err := useCase.DeleteData("1")
	assert.Error(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestFindByIdSuccess() {
	suite.repoMock.On("Get", "1").Return(&brandDummies[0], nil)
	useCase := NewBrandUseCase(suite.repoMock)
	brand, err := useCase.FindById("1")
	assert.Equal(suite.T(), brandDummies[0], *brand)
	assert.Nil(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestFindByIdRepoErrorFail() {
	suite.repoMock.On("Get", "1").Return(nil, errors.New(repositoryErrorMessage))
	useCase := NewBrandUseCase(suite.repoMock)
	brand, err := useCase.FindById("1")
	assert.Nil(suite.T(), brand)
	assert.Error(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestSearchBySuccess() {
	filter := map[string]interface{}{"brand": "Honda"}
	suite.repoMock.On("Search", filter).Return(brandDummies, nil)
	useCase := NewBrandUseCase(suite.repoMock)
	brands, err := useCase.SearchBy(filter)
	assert.Equal(suite.T(), brandDummies, brands)
	assert.Nil(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestSearchByRepoErrorFail() {
	filter := map[string]interface{}{"brand": "Honda"}
	suite.repoMock.On("Search", filter).Return(nil, errors.New(repositoryErrorMessage))
	useCase := NewBrandUseCase(suite.repoMock)
	brands, err := useCase.SearchBy(filter)
	assert.Nil(suite.T(), brands)
	assert.Error(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestSaveDataSuccess() {
	dummy := brandDummies[0]
	var countBrand int64 = 0
	suite.repoMock.On("CountByName", "Honda", "1").Return(countBrand, nil)
	suite.repoMock.On("Get", "1").Return(&brandDummies[0], nil)
	suite.repoMock.On("Save", &dummy).Return(nil)
	useCase := NewBrandUseCase(suite.repoMock)
	err := useCase.SaveData(&dummy)
	assert.Nil(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestSaveDataRepoErrorFail() {
	dummy := brandDummies[0]
	var countBrand int64 = 1
	suite.repoMock.On("CountByName", "Honda", "1").Return(countBrand, errors.New(repositoryErrorMessage))
	suite.repoMock.On("Save", &dummy).Return(errors.New(repositoryErrorMessage))
	useCase := NewBrandUseCase(suite.repoMock)
	err := dummy.Validate()
	assert.Error(suite.T(), err)
	dummy.Name = ""
	err = useCase.SaveData(&dummy)
	assert.Error(suite.T(), err)
	dummy = brandDummies[0]
	err = useCase.SaveData(&dummy)
	assert.Error(suite.T(), err)
}

func (suite *BrandUseCaseTestSuite) TestSaveDataIDNotFoundFail() {
	dummy := brandDummies[0]
	var countBrand int64 = 0
	suite.repoMock.On("CountByName", "Honda", "1").Return(countBrand, nil)
	suite.repoMock.On("Get", "1").Return(nil, errors.New("not found"))
	useCase := NewBrandUseCase(suite.repoMock)
	err := useCase.SaveData(&dummy)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "brand with ID 1 not found", err.Error())
}

func (suite *BrandUseCaseTestSuite) TestPaginationSuccess() {
	brandDm := brandDummies
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   5,
		TotalPages:  1,
	}
	suite.repoMock.On("Paging", mock.AnythingOfType("dto.RequestQueryParams")).Return(brandDm, expectedPaging, nil)
	useCase := NewBrandUseCase(suite.repoMock)
	requestParams := dto.RequestQueryParams{QueryParams: dto.QueryParams{Sort: "ASC"}}
	actualBrand, actualPaging, actualError := useCase.Pagination(requestParams)
	assert.Equal(suite.T(), brandDm, actualBrand)
	assert.Equal(suite.T(), expectedPaging, actualPaging)
	assert.Equal(suite.T(), nil, actualError)
}

func (suite *BrandUseCaseTestSuite) TestPaginationFail() {
	expectedPaging := dto.Paging{
		Page:        0,
		RowsPerPage: 0,
		TotalRows:   0,
		TotalPages:  0,
	}
	suite.repoMock.On("Paging", mock.AnythingOfType("dto.RequestQueryParams")).Return(nil, expectedPaging, errors.New(repositoryErrorMessage))
	useCase := NewBrandUseCase(suite.repoMock)
	requestParams := dto.RequestQueryParams{QueryParams: dto.QueryParams{Sort: "ABC"}}
	_, actualPaging, actualError := useCase.Pagination(requestParams)
	assert.Equal(suite.T(), expectedPaging, actualPaging)
	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), "invalid sort by: ABC", actualError.Error())
}

// BrandUseCaseTestSuite as test suite model, any field suite and repoMock
type BrandUseCaseTestSuite struct {
	suite.Suite
	repoMock *repoMock
}

// BrandUseCaseTestSuite as SetupTest from repoMock
func (suite *BrandUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

// TestBrandUseCaseTestSuite as constructor BrandUseCase and running all  test
func TestBrandUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(BrandUseCaseTestSuite))
}
