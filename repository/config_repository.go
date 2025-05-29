package repository

import (
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/gorm"
)

type ConfigRepository interface {
	Create(config *models.Config) error
	GetAll() ([]models.Config, error)
	GetByKey(key string) (*models.Config, error)
	Update(config *models.Config) error
	Delete(key string) error
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) Create(config *models.Config) error {
	return r.db.Create(config).Error
}

func (r *configRepository) GetAll() ([]models.Config, error) {
	var configs []models.Config
	err := r.db.Find(&configs).Error
	return configs, err
}

func (r *configRepository) GetByKey(key string) (*models.Config, error) {
	var config models.Config
	err := r.db.First(&config, "key = ?", key).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *configRepository) Update(config *models.Config) error {
	return r.db.Save(config).Error
}

func (r *configRepository) Delete(key string) error {
	return r.db.Delete(&models.Config{}, "key = ?", key).Error
}

