package services

import (
	"assets-api-go/internal/common"
	"assets-api-go/internal/dto"
	"assets-api-go/internal/models"
	"assets-api-go/internal/repositories"
	"log"
	"net/http"
	"time"
)

type AssetServiceInterface interface {
	CreateAsset(input *dto.AssetInputDto) (code int, response *dto.BaseResponse)
	GetAssetById(id string) (code int, response *dto.BaseResponse)
	GetAssets(pagination *dto.MetaPagination) (code int, response *dto.MetaPagination)
	UpdateAsset(id string, input *dto.AssetInputDto) (code int, response *dto.BaseResponse)
	DeleteAsset(id string) (code int, response *dto.BaseResponse)
}

type assetService struct {
	assetRepo repositories.AssetRepositoryInterface
}

func NewAssetService(assetRepo repositories.AssetRepositoryInterface) AssetServiceInterface {
	return &assetService{assetRepo: assetRepo}
}

func (s *assetService) CreateAsset(input *dto.AssetInputDto) (code int, response *dto.BaseResponse) {
	// get asset by name and type
	asset, err := s.assetRepo.GetAssetByAttribute(map[string]interface{}{
		"name": input.Name,
		"type": input.Type,
	})
	if err != nil {
		log.Println("[assetService][CreateAsset] error get existing asset :", err)
		return http.StatusInternalServerError, &dto.BaseResponse{
			Error:            common.InternalServerError,
			ErrorDescription: "Something went wrong",
		}
	}

	if asset != nil {
		return http.StatusBadRequest, &dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "Asset already exist",
		}
	}

	acqusitionDate, err := time.Parse("2006-01-02", input.AcquisitionDate)
	if err != nil {
		log.Println("[assetService][CreateAsset] error parsing date :", err)
		return http.StatusBadRequest, &dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "Invalid acqusition date format",
		}
	}
	tx := s.assetRepo.StartTransaction()
	asset, err = s.assetRepo.CreateAsset(&models.Asset{
		Name:            input.Name,
		Type:            input.Type,
		Value:           input.Value,
		AcquisitionDate: acqusitionDate,
	}, tx)
	if err != nil {
		log.Println("[assetService][CreateAsset] error create asset :", err)
		err = s.assetRepo.RollbackTransaction(tx)
		if err != nil {
			log.Println("[assetService][CreateAsset] error rollback transaction :", err)
		}
		return http.StatusInternalServerError, &dto.BaseResponse{
			Error:            common.InternalServerError,
			ErrorDescription: "Something went wrong",
		}
	}

	resData := &dto.AssetOutputDto{
		Id:              asset.Id,
		Name:            asset.Name,
		Type:            asset.Type,
		Value:           asset.Value,
		AcquisitionDate: asset.AcquisitionDate.Format("2006-01-02 15:04:05"),
		CreatedAt:       asset.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       asset.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return http.StatusCreated, &dto.BaseResponse{
		Message: common.Success,
		Data:    resData,
	}
}

func (s *assetService) GetAssetById(id string) (code int, response *dto.BaseResponse) {
	// get asset by name and type
	asset, err := s.assetRepo.GetAssetByAttribute(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		log.Println("[assetService][GetAssetById] error get existing asset :", err)
		return http.StatusInternalServerError, &dto.BaseResponse{
			Error:            common.InternalServerError,
			ErrorDescription: "Something went wrong",
		}
	}

	if asset == nil {
		return http.StatusNotFound, &dto.BaseResponse{
			Error:            common.NotFound,
			ErrorDescription: "Asset not found",
		}
	}

	resData := &dto.AssetOutputDto{
		Id:              asset.Id,
		Name:            asset.Name,
		Type:            asset.Type,
		Value:           asset.Value,
		AcquisitionDate: asset.AcquisitionDate.Format("2006-01-02 15:04:05"),
		CreatedAt:       asset.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       asset.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return http.StatusOK, &dto.BaseResponse{
		Message: common.Success,
		Data:    resData,
	}
}

func (s *assetService) GetAssets(pagination *dto.MetaPagination) (code int, response *dto.MetaPagination) {

	assets, count, err := s.assetRepo.GetAssets(pagination)
	if err != nil {
		log.Println("[assetService][GetAssets] error get assets :", err)
		return http.StatusInternalServerError, &dto.MetaPagination{
			BaseResponse: dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		}
	}

	assetsRes := []*dto.AssetOutputDto{}
	for _, v := range assets {
		assetsRes = append(assetsRes, &dto.AssetOutputDto{
			Id:              v.Id,
			Name:            v.Name,
			Type:            v.Type,
			Value:           v.Value,
			AcquisitionDate: v.AcquisitionDate.Format("2006-01-02"),
			CreatedAt:       v.CreatedAt.Format("2006-01-02"),
			UpdatedAt:       v.UpdatedAt.Format("2006-01-02"),
		})
	}
	pagination.Total = count
	pagination.TotalPage = count / int64(pagination.Limit)
	if count%int64(pagination.Limit) > 0 {
		pagination.TotalPage++
	}
	pagination.Data = assetsRes
	return http.StatusOK, pagination
}

func (s *assetService) UpdateAsset(id string, input *dto.AssetInputDto) (code int, response *dto.BaseResponse) {
	asset, err := s.assetRepo.GetAssetByAttribute(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		log.Println("[assetService][UpdateAsset] error get existing asset :", err)
		return http.StatusInternalServerError, &dto.BaseResponse{
			Error:            common.InternalServerError,
			ErrorDescription: "Something went wrong",
		}
	}

	if asset == nil {
		return http.StatusNotFound, &dto.BaseResponse{
			Error:            common.NotFound,
			ErrorDescription: "Asset not found",
		}
	}

	acqusitionDate, err := time.Parse("2006-01-02", input.AcquisitionDate)
	if err != nil {
		log.Println("[assetService][UpdateAsset] error parsing date :", err)
		return http.StatusBadRequest, &dto.BaseResponse{
			Error:            common.BadRequest,
			ErrorDescription: "Invalid acqusition date format",
		}
	}

	asset.Name = input.Name
	asset.Type = input.Type
	asset.Value = input.Value
	asset.AcquisitionDate = acqusitionDate

	tx := s.assetRepo.StartTransaction()
	asset, err = s.assetRepo.UpdateAsset(asset, tx)
	if err != nil {
		log.Println("[assetService][UpdateAsset] error update asset :", err)
		err = s.assetRepo.RollbackTransaction(tx)
		if err != nil {
			log.Println("[assetService][UpdateAsset] error rollback transaction :", err)
		}
		return http.StatusInternalServerError, &dto.BaseResponse{
			Error:            common.InternalServerError,
			ErrorDescription: "Something went wrong",
		}
	}

	resData := &dto.AssetOutputDto{
		Id:              asset.Id,
		Name:            asset.Name,
		Type:            asset.Type,
		Value:           asset.Value,
		AcquisitionDate: asset.AcquisitionDate.Format("2006-01-02 15:04:05"),
		CreatedAt:       asset.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       asset.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return http.StatusOK, &dto.BaseResponse{
		Message: common.Success,
		Data:    resData,
	}
}

func (s *assetService) DeleteAsset(id string) (code int, response *dto.BaseResponse) {

	asset, err := s.assetRepo.GetAssetByAttribute(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		log.Println("[assetService][UpdateAsset] error get existing asset :", err)
		return http.StatusInternalServerError, &dto.BaseResponse{
			Error:            common.InternalServerError,
			ErrorDescription: "Something went wrong",
		}
	}

	if asset == nil {
		return http.StatusNotFound, &dto.BaseResponse{
			Error:            common.NotFound,
			ErrorDescription: "Asset not found",
		}
	}

	tx := s.assetRepo.StartTransaction()
	err = s.assetRepo.DeleteAsset(asset, tx)
	if err != nil {
		log.Println("[assetService][DeleteAsset] error delete asset :", err)
		err = s.assetRepo.RollbackTransaction(tx)
		if err != nil {
			log.Println("[assetService][DeleteAsset] error rollback transaction :", err)
		}
		return http.StatusInternalServerError, &dto.BaseResponse{
			Error:            common.InternalServerError,
			ErrorDescription: "Something went wrong",
		}
	}

	return http.StatusOK, &dto.BaseResponse{
		Message: common.Success,
	}
}
