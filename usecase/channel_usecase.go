package usecase

import (
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
)

type ChannelUsecase interface {
	CreateChannel(channel *models.Channel) error
	GetAllChannels() ([]models.Channel, error)
	GetChannelByID(id string) (*models.Channel, error)
	UpdateChannel(channel *models.Channel) error
	DeleteChannel(id string) error
}

type channelUsecase struct {
	channelRepo repository.ChannelRepository
}

func NewChannelUsecase(channelRepo repository.ChannelRepository) ChannelUsecase {
	return &channelUsecase{channelRepo: channelRepo}
}

func (u *channelUsecase) CreateChannel(channel *models.Channel) error {
	return u.channelRepo.Create(channel)
}

func (u *channelUsecase) GetAllChannels() ([]models.Channel, error) {
	return u.channelRepo.GetAll()
}

func (u *channelUsecase) GetChannelByID(id string) (*models.Channel, error) {
	return u.channelRepo.GetByID(id)
}

func (u *channelUsecase) UpdateChannel(channel *models.Channel) error {
	return u.channelRepo.Update(channel)
}

func (u *channelUsecase) DeleteChannel(id string) error {
	return u.channelRepo.Delete(id)
}
