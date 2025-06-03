package usecase

import (
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
	"github.com/kougami132/MsgPilot/internal/types"
)

type MessageUsecase interface {
	CreateMessage(message *models.Message) error
	GetAllMessages() ([]models.Message, error)
	GetMessageByID(id int) (*models.Message, error)
	UpdateMessage(message *models.Message) error
	UpdateMessageStatus(message *models.Message, status types.MessageStatus) error
	UpdateMessageStatusWithErrorMessage(message *models.Message, status types.MessageStatus, errorMessage string) error
	DeleteMessage(id int) error
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

func (u *messageUsecase) GetMessageByID(id int) (*models.Message, error) {
	return u.messageRepo.GetByID(id)
}

func (u *messageUsecase) UpdateMessage(message *models.Message) error {
	return u.messageRepo.Update(message)
}

func (u *messageUsecase) UpdateMessageStatus(message *models.Message, status types.MessageStatus) error {
	message.Status = status
	return u.UpdateMessage(message)
}

func (u *messageUsecase) UpdateMessageStatusWithErrorMessage(message *models.Message, status types.MessageStatus, errorMessage string) error {
	message.Status = status
	message.ErrorMessage = errorMessage
	return u.UpdateMessage(message)
}

func (u *messageUsecase) DeleteMessage(id int) error {
	return u.messageRepo.Delete(id)
}
