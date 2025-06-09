package usecase

import (
	"github.com/kougami132/MsgPilot/channels"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
)

type ChannelUsecase interface {
	CreateChannel(channel *models.Channel) error
	GetAllChannels() ([]models.Channel, error)
	GetChannelByID(id int) (*models.Channel, error)
	UpdateChannel(channel *models.Channel) error
	DeleteChannel(id int) error
	TestPush(channel models.Channel) error
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

func (u *channelUsecase) GetChannelByID(id int) (*models.Channel, error) {
	return u.channelRepo.GetByID(id)
}

func (u *channelUsecase) UpdateChannel(channel *models.Channel) error {
	return u.channelRepo.Update(channel)
}

func (u *channelUsecase) DeleteChannel(id int) error {
	return u.channelRepo.Delete(id)
}


func (u *channelUsecase) TestPush(channel models.Channel) error {
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