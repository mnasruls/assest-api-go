package services

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"assets-api-go/internal/common"
	"assets-api-go/internal/dto"
	"assets-api-go/internal/models"
	"assets-api-go/mocks/repositories"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAssetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockAssetRepositoryInterface(ctrl)
	service := NewAssetService(mockRepo)

	testTime := time.Now()
	testAsset := &models.Asset{
		Id:              "test-id",
		Name:            "Test Asset",
		Type:            "Test Type",
		Value:           1000,
		AcquisitionDate: testTime,
		CreatedAt:       testTime,
		UpdatedAt:       testTime,
	}

	tests := []struct {
		name           string
		id             string
		mockSetup      func()
		expectedCode   int
		expectedResult *dto.BaseResponse
	}{
		{
			name: "Success - Asset found",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
			},
			expectedCode: http.StatusOK,
			expectedResult: &dto.BaseResponse{
				Message: common.Success,
				Data: &dto.AssetOutputDto{
					Id:              "test-id",
					Name:            "Test Asset",
					Type:            "Test Type",
					Value:           1000,
					AcquisitionDate: testTime.Format("2006-01-02 15:04:05"),
					CreatedAt:       testTime.Format("2006-01-02 15:04:05"),
					UpdatedAt:       testTime.Format("2006-01-02 15:04:05"),
				},
			},
		},
		{
			name: "Error - Asset not found",
			id:   "non-existent-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "non-existent-id"}).Return(nil, nil)
			},
			expectedCode: http.StatusNotFound,
			expectedResult: &dto.BaseResponse{
				Error:            common.NotFound,
				ErrorDescription: "Asset not found",
			},
		},
		{
			name: "Error - Repository error",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(nil, errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			code, response := service.GetAssetById(tt.id)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedResult, response)
		})
	}
}

func TestCreateAsset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockAssetRepositoryInterface(ctrl)
	service := NewAssetService(mockRepo)

	tests := []struct {
		name           string
		input          *dto.AssetInputDto
		mockSetup      func()
		expectedCode   int
		expectedResult *dto.BaseResponse
	}{
		{
			name: "Success - Create new asset",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(nil, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().CreateAsset(gomock.Any(), gomock.Any()).Return(&models.Asset{
					Id:              "test-id",
					Name:            "Test Asset",
					Type:            "Test Type",
					Value:           1000,
					AcquisitionDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedResult: &dto.BaseResponse{
				Message: common.Success,
				Data: &dto.AssetOutputDto{
					Id:              "test-id",
					Name:            "Test Asset",
					Type:            "Test Type",
					Value:           1000,
					AcquisitionDate: "2023-01-01 00:00:00",
				},
			},
		},
		{
			name: "Error - Asset already exists",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(&models.Asset{}, nil)
			},
			expectedCode: http.StatusBadRequest,
			expectedResult: &dto.BaseResponse{
				Error:            common.BadRequest,
				ErrorDescription: "Asset already exist",
			},
		},
		{
			name: "Error - Invalid date format",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "invalid-date",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(nil, nil)
			},
			expectedCode: http.StatusBadRequest,
			expectedResult: &dto.BaseResponse{
				Error:            common.BadRequest,
				ErrorDescription: "Invalid acqusition date format",
			},
		},
		{
			name: "Error - Get Exist asset error",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(nil, errors.New("get exist asset error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Failed to create asset",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(nil, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().CreateAsset(gomock.Any(), gomock.Any()).Return(nil, errors.New("create error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Failed to commit transaction",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(nil, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().CreateAsset(gomock.Any(), gomock.Any()).Return(&models.Asset{}, nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(errors.New("commit error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Failed to rollback transaction 1",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(nil, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().CreateAsset(gomock.Any(), gomock.Any()).Return(&models.Asset{}, nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(errors.New("commit error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(errors.New("rollback error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Failed to rollback transaction 2",
			input: &dto.AssetInputDto{
				Name:            "Test Asset",
				Type:            "Test Type",
				Value:           1000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(gomock.Any()).Return(nil, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().CreateAsset(gomock.Any(), gomock.Any()).Return(nil, errors.New("create error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(errors.New("rollback error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			code, _ := service.CreateAsset(tt.input)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

func TestGetAssets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockAssetRepositoryInterface(ctrl)
	service := NewAssetService(mockRepo)

	testTime := time.Now()
	testAssets := []*models.Asset{
		{
			Id:              "test-id-1",
			Name:            "Test Asset 1",
			Type:            "Test Type",
			Value:           1000,
			AcquisitionDate: testTime,
			CreatedAt:       testTime,
			UpdatedAt:       testTime,
		},
		{
			Id:              "test-id-2",
			Name:            "Test Asset 2",
			Type:            "Test Type",
			Value:           2000,
			AcquisitionDate: testTime,
			CreatedAt:       testTime,
			UpdatedAt:       testTime,
		},
	}

	tests := []struct {
		name           string
		pagination     *dto.MetaPagination
		mockSetup      func()
		expectedCode   int
		expectedResult *dto.MetaPagination
	}{
		{
			name: "Success - Get assets with pagination",
			pagination: &dto.MetaPagination{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssets(gomock.Any()).Return(testAssets, int64(2), nil)
			},
			expectedCode: http.StatusOK,
			expectedResult: &dto.MetaPagination{
				Limit:     10,
				Offset:    0,
				Total:     2,
				TotalPage: 1,
				BaseResponse: dto.BaseResponse{
					Data: []*dto.AssetOutputDto{
						{
							Id:              "test-id-1",
							Name:            "Test Asset 1",
							Type:            "Test Type",
							Value:           1000,
							AcquisitionDate: testTime.Format("2006-01-02"),
							CreatedAt:       testTime.Format("2006-01-02"),
							UpdatedAt:       testTime.Format("2006-01-02"),
						},
						{
							Id:              "test-id-2",
							Name:            "Test Asset 2",
							Type:            "Test Type",
							Value:           2000,
							AcquisitionDate: testTime.Format("2006-01-02"),
							CreatedAt:       testTime.Format("2006-01-02"),
							UpdatedAt:       testTime.Format("2006-01-02"),
						},
					},
				},
			},
		},
		{
			name: "Error - Repository error",
			pagination: &dto.MetaPagination{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssets(gomock.Any()).Return(nil, int64(0), errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.MetaPagination{
				BaseResponse: dto.BaseResponse{
					Error:            common.InternalServerError,
					ErrorDescription: "Something went wrong",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			code, response := service.GetAssets(tt.pagination)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedResult, response)
		})
	}
}

func TestUpdateAsset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockAssetRepositoryInterface(ctrl)
	service := NewAssetService(mockRepo)

	testTime := time.Now()
	testAsset := &models.Asset{
		Id:              "test-id",
		Name:            "Test Asset",
		Type:            "Test Type",
		Value:           1000,
		AcquisitionDate: testTime,
		CreatedAt:       testTime,
		UpdatedAt:       testTime,
	}

	tests := []struct {
		name           string
		id             string
		input          *dto.AssetInputDto
		mockSetup      func()
		expectedCode   int
		expectedResult *dto.BaseResponse
	}{
		{
			name: "Success - Update asset",
			id:   "test-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().UpdateAsset(gomock.Any(), gomock.Any()).Return(testAsset, nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedResult: &dto.BaseResponse{
				Message: common.Success,
				Data: &dto.AssetOutputDto{
					Id:              "test-id",
					Name:            "Test Asset",
					Type:            "Test Type",
					Value:           1000,
					AcquisitionDate: testTime.Format("2006-01-02 15:04:05"),
					CreatedAt:       testTime.Format("2006-01-02 15:04:05"),
					UpdatedAt:       testTime.Format("2006-01-02 15:04:05"),
				},
			},
		},
		{
			name: "Error - Asset not found",
			id:   "non-existent-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "non-existent-id"}).Return(nil, nil)
			},
			expectedCode: http.StatusNotFound,
			expectedResult: &dto.BaseResponse{
				Error:            common.NotFound,
				ErrorDescription: "Asset not found",
			},
		},
		{
			name: "Error - Invalid date format",
			id:   "test-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "invalid-date",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
			},
			expectedCode: http.StatusBadRequest,
			expectedResult: &dto.BaseResponse{
				Error:            common.BadRequest,
				ErrorDescription: "Invalid acqusition date format",
			},
		},
		{
			name: "Error - Get asset error",
			id:   "test-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(nil, errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Update asset error",
			id:   "test-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().UpdateAsset(gomock.Any(), gomock.Any()).Return(nil, errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Commit transaction error",
			id:   "test-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().UpdateAsset(gomock.Any(), gomock.Any()).Return(testAsset, nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Rollback transaction error 1",
			id:   "test-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().UpdateAsset(gomock.Any(), gomock.Any()).Return(testAsset, nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Rollback transaction error 2",
			id:   "test-id",
			input: &dto.AssetInputDto{
				Name:            "Updated Asset",
				Type:            "Updated Type",
				Value:           2000,
				AcquisitionDate: "2023-01-01",
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().UpdateAsset(gomock.Any(), gomock.Any()).Return(testAsset, errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			code, _ := service.UpdateAsset(tt.id, tt.input)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

func TestDeleteAsset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockAssetRepositoryInterface(ctrl)
	service := NewAssetService(mockRepo)

	testAsset := &models.Asset{
		Id: "test-id",
	}

	tests := []struct {
		name           string
		id             string
		mockSetup      func()
		expectedCode   int
		expectedResult *dto.BaseResponse
	}{
		{
			name: "Success - Delete asset",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().DeleteAsset(testAsset, gomock.Any()).Return(nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedResult: &dto.BaseResponse{
				Message: common.Success,
			},
		},
		{
			name: "Error - Asset not found",
			id:   "non-existent-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "non-existent-id"}).Return(nil, nil)
			},
			expectedCode: http.StatusNotFound,
			expectedResult: &dto.BaseResponse{
				Error:            common.NotFound,
				ErrorDescription: "Asset not found",
			},
		},
		{
			name: "Error - Repository error on get",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(nil, errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Repository error on delete",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().DeleteAsset(testAsset, gomock.Any()).Return(errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Repository error on commit",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().DeleteAsset(testAsset, gomock.Any()).Return(nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Repository error on rollback 1",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().DeleteAsset(testAsset, gomock.Any()).Return(errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
		{
			name: "Error - Repository error on rollback 2",
			id:   "test-id",
			mockSetup: func() {
				mockRepo.EXPECT().GetAssetByAttribute(map[string]interface{}{"id": "test-id"}).Return(testAsset, nil)
				mockRepo.EXPECT().StartTransaction().Return(&gorm.DB{})
				mockRepo.EXPECT().DeleteAsset(testAsset, gomock.Any()).Return(nil)
				mockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(errors.New("repository error"))
				mockRepo.EXPECT().RollbackTransaction(gomock.Any()).Return(errors.New("repository error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResult: &dto.BaseResponse{
				Error:            common.InternalServerError,
				ErrorDescription: "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			code, response := service.DeleteAsset(tt.id)
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedResult, response)
		})
	}
}
