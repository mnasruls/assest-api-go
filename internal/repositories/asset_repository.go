package repositories

import (
	"assets-api-go/internal/dto"
	"assets-api-go/internal/models"

	"gorm.io/gorm"
)

type AssetRepositoryInterface interface {
	StartTransaction() *gorm.DB
	CommitTransaction(*gorm.DB) error
	RollbackTransaction(*gorm.DB) error
	CreateAsset(asset *models.Asset, tx *gorm.DB) (*models.Asset, error)
	GetAssetByAttribute(whereClause interface{}) (*models.Asset, error)
	GetAssets(pagination *dto.MetaPagination) ([]*models.Asset, int64, error)
	UpdateAsset(asset *models.Asset, tx *gorm.DB) (*models.Asset, error)
	DeleteAsset(asset *models.Asset, tx *gorm.DB) error
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepositoryInterface {
	return &assetRepository{db}
}

func (repo *assetRepository) StartTransaction() *gorm.DB {
	return repo.db.Begin()
}

func (repo *assetRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (repo *assetRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *assetRepository) CreateAsset(asset *models.Asset, tx *gorm.DB) (*models.Asset, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Create(asset).Error; err != nil {
		return nil, err
	}

	return asset, nil
}

func (r *assetRepository) GetAssetByAttribute(whereClause interface{}) (*models.Asset, error) {
	var asset models.Asset

	if err := r.db.Where(whereClause).Where("deleted_at is NULL").Order("created_at desc").First(&asset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &asset, nil
}

func (r *assetRepository) GetAssets(pagination *dto.MetaPagination) ([]*models.Asset, int64, error) {
	var assets []*models.Asset
	var total int64

	query := r.db.Where("deleted_at is NULL").Order("created_at desc")

	if err := query.Model(&models.Asset{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(pagination.Limit).Offset(pagination.Offset).Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
}

func (r *assetRepository) UpdateAsset(asset *models.Asset, tx *gorm.DB) (*models.Asset, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.Save(asset).Error; err != nil {
		return nil, err
	}

	return asset, nil
}

func (r *assetRepository) DeleteAsset(asset *models.Asset, tx *gorm.DB) error {

	if err := r.db.Delete(&asset).Error; err != nil {
		return err
	}

	return nil
}
