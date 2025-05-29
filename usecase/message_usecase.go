package usecase

import (
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
)

type MessageUsecase interface {
	CreateMessage(message *models.Message) error
	GetAllMessages() ([]models.Message, error)
	GetMessageByID(id string) (*models.Message, error)
	UpdateMessage(message *models.Message) error
	DeleteMessage(id string) error
}

type messageUsecase struct {
	messageRepo repository.MessageRepository
}

func NewMessageUsecase(messageRepo repository.MessageRepository) MessageUsecase {
	return &messageUsecase{messageRepo: messageRepo}
}

func (u *messageUsecase) CreateMessage(message *models.Message) error {
	return u.messageRepo.Create(message)
}

func (u *messageUsecase) GetAllMessages() ([]models.Message, error) {
	return u.messageRepo.GetAll()
}

func (u *messageUsecase) GetMessageByID(id string) (*models.Message, error) {
	return u.messageRepo.GetByID(id)
}

func (u *messageUsecase) UpdateMessage(message *models.Message) error {
	return u.messageRepo.Update(message)
}

func (u *messageUsecase) DeleteMessage(id string) error {
	return u.messageRepo.Delete(id)
}
