package service

import (
	"github.com/kougami132/MsgPilot/internal/repository"
	"github.com/kougami132/MsgPilot/internal/channels"
	"github.com/kougami132/MsgPilot/models"
)

type ChannelService interface {
	CreateChannel(channel *models.Channel) error
	GetAllChannels() ([]models.Channel, error)
	GetChannelByID(id int) (*models.Channel, error)
	UpdateChannel(channel *models.Channel) error
	DeleteChannel(id int) error
	TestPush(channel models.Channel) error
}

type channelService struct {
	channelRepo repository.ChannelRepository
}

func NewChannelService(channelRepo repository.ChannelRepository) ChannelService {
	return &channelService{channelRepo: channelRepo}
}

func (u *channelService) CreateChannel(channel *models.Channel) error {
	return u.channelRepo.Create(channel)
}

func (u *channelService) GetAllChannels() ([]models.Channel, error) {
	return u.channelRepo.GetAll()
}

func (u *channelService) GetChannelByID(id int) (*models.Channel, error) {
	return u.channelRepo.GetByID(id)
}

func (u *channelService) UpdateChannel(channel *models.Channel) error {
	return u.channelRepo.Update(channel)
}

func (u *channelService) DeleteChannel(id int) error {
	return u.channelRepo.Delete(id)
}


func (u *channelService) TestPush(channel models.Channel) error {
	testMessage := &models.Message{
		Title:    "MsgPilot消息推送",
		Content:  "测试消息",
	}

	// 发送消息
	handler, err := channels.GetChannelHandler(channel)
	if err != nil {
		return err
	}
	err = handler.Send(testMessage)
	if err != nil {
		return err
	}

	return nil
}