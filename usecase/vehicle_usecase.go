package usecase

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/repository"
)

type VehicleUseCase interface {
	BaseUseCase[model.Vehicle]
	Paging(requestQueryParams dto.RequestQueryParams) ([]model.Vehicle, dto.Paging, error)
	UpdateVehicleStock(count int, id string) error
	UploadImage(payload *model.Vehicle, file multipart.File, fileExt string) error
}

type vehicleUseCase struct {
	repo         repository.VehicleRepository
	brandUseCase BrandUseCase
	fileUseCase  FileUseCase
}

func (v *vehicleUseCase) SearchBy(by map[string]interface{}) ([]model.Vehicle, error) {
	return v.repo.Search(by)
}

func (v *vehicleUseCase) FindAll() ([]model.Vehicle, error) {
	return v.repo.List()
}

func (v *vehicleUseCase) FindById(id string) (*model.Vehicle, error) {
	vehicle, err := v.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("vehicle with id %s not found", id)
	}

	return vehicle, nil
}

func (v *vehicleUseCase) SaveData(payload *model.Vehicle) error {
	brand, err := v.brandUseCase.FindById(payload.BrandID)
	if err != nil {
		return fmt.Errorf("Brand with ID %s not found!", payload.ID)
	}
	payload.BrandID = brand.ID
	return v.repo.Save(payload)
}

func (v *vehicleUseCase) DeleteData(id string) error {
	return v.repo.Delete(id)
}

func (v *vehicleUseCase) UpdateVehicleStock(count int, id string) error {
	return v.repo.UpdateStock(count, id)
}

func (v *vehicleUseCase) Paging(requestQueryParams dto.RequestQueryParams) ([]model.Vehicle, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return v.repo.Paging(requestQueryParams)

}

func (v *vehicleUseCase) UploadImage(payload *model.Vehicle, file multipart.File, fileExt string) error {
	concatName := fmt.Sprintf("%s-%s", payload.Model, payload.BrandID)
	fileName := fmt.Sprintf("img-%s.%s", strings.ToLower(concatName), fileExt)
	fileLocation, err := v.fileUseCase.Save(file, fileName)
	if err != nil {
		return err
	}

	payload.ImgPath = fileLocation
	payload.UrlPath = fmt.Sprintf("/vehicles/image/%s", strings.ToLower(concatName))
	err = v.SaveData(payload)
	if err != nil {
		return err
	}
	return nil
}

func NewVehicleUseCase(
	repo repository.VehicleRepository,
	brandUseCase BrandUseCase,
	fileUseCase FileUseCase,
) VehicleUseCase {
	return &vehicleUseCase{
		repo:         repo,
		brandUseCase: brandUseCase,
		fileUseCase:  fileUseCase,
	}
}
