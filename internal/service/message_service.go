package service

import (
	"github.com/kougami132/MsgPilot/internal/repository"
	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
)

type MessageService interface {
	CreateMessage(message *models.Message) error
	GetAllMessages() ([]models.Message, error)
	GetMessageByID(id int) (*models.Message, error)
	UpdateMessage(message *models.Message) error
	UpdateMessageStatus(message *models.Message, status types.MessageStatus) error
	UpdateMessageStatusWithErrorMessage(message *models.Message, status types.MessageStatus, errorMessage string) error
	DeleteMessage(id int) error
}

type messageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(messageRepo repository.MessageRepository) MessageService {
	return &messageService{messageRepo: messageRepo}
}

func (u *messageService) CreateMessage(message *models.Message) error {
	return u.messageRepo.Create(message)
}

func (u *messageService) GetAllMessages() ([]models.Message, error) {
	return u.messageRepo.GetAll()
}

func (u *messageService) GetMessageByID(id int) (*models.Message, error) {
	return u.messageRepo.GetByID(id)
}

func (u *messageService) UpdateMessage(message *models.Message) error {
	return u.messageRepo.Update(message)
}

func (u *messageService) UpdateMessageStatus(message *models.Message, status types.MessageStatus) error {
	message.Status = status
	return u.UpdateMessage(message)
}

func (u *messageService) UpdateMessageStatusWithErrorMessage(message *models.Message, status types.MessageStatus, errorMessage string) error {
	message.Status = status
	message.ErrorMessage = errorMessage
	return u.UpdateMessage(message)
}

func (u *messageService) DeleteMessage(id int) error {
	return u.messageRepo.Delete(id)
}
