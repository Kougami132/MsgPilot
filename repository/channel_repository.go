package repository

import (
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/gorm"
)

type ChannelRepository interface {
	Create(channel *models.Channel) error
	GetAll() ([]models.Channel, error)
	GetByID(id int) (*models.Channel, error)
	Update(channel *models.Channel) error
	Delete(id int) error
}

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) Create(channel *models.Channel) error {
	return r.db.Create(channel).Error
}

func (r *channelRepository) GetAll() ([]models.Channel, error) {
	var channels []models.Channel
	err := r.db.Find(&channels).Error
	return channels, err
}

func (r *channelRepository) GetByID(id int) (*models.Channel, error) {
	var channel models.Channel
	err := r.db.First(&channel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

func (r *channelRepository) Update(channel *models.Channel) error {
	return r.db.Save(channel).Error
}

func (r *channelRepository) Delete(id int) error {
	return r.db.Delete(&models.Channel{}, "id = ?", id).Error
}
