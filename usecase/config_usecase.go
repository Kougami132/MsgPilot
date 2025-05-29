package usecase

import (
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
)

type ConfigUsecase interface {
	CreateConfig(config *models.Config) error
	GetAllConfigs() ([]models.Config, error)
	GetConfigByKey(key string) (*models.Config, error)
	UpdateConfig(config *models.Config) error
	DeleteConfig(key string) error
	InitConfig(key string, defaultValue string) error
}

type configUsecase struct {
	configRepo repository.ConfigRepository
}

func NewConfigUsecase(configRepo repository.ConfigRepository) ConfigUsecase {
	return &configUsecase{configRepo: configRepo}
}

func (u *configUsecase) CreateConfig(config *models.Config) error {
	return u.configRepo.Create(config)
}

func (u *configUsecase) GetAllConfigs() ([]models.Config, error) {
	return u.configRepo.GetAll()
}

func (u *configUsecase) GetConfigByKey(key string) (*models.Config, error) {
	return u.configRepo.GetByKey(key)
}

func (u *configUsecase) UpdateConfig(config *models.Config) error {
	return u.configRepo.Update(config)
}

func (u *configUsecase) DeleteConfig(key string) error {
	return u.configRepo.Delete(key)
}

func (u *configUsecase) InitConfig(key string, defaultValue string) error {
	if _, err := u.GetConfigByKey(key); err != nil {
		return u.CreateConfig(&models.Config{Key: key, Value: defaultValue})
	}
	return nil
}
