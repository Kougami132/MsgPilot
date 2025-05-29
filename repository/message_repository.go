package repository

import (
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *models.Message) error
	GetAll() ([]models.Message, error)
	GetByID(id string) (*models.Message, error)
	Update(message *models.Message) error
	Delete(id string) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetAll() ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Preload("ReceiveChannel").Preload("SendChannel").Find(&messages).Error
	return messages, err
}

func (r *messageRepository) GetByID(id string) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("ReceiveChannel").Preload("SendChannel").First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *messageRepository) Delete(id string) error {
	return r.db.Delete(&models.Message{}, "id = ?", id).Error
}
