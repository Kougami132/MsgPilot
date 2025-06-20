package service

import (
	"errors"

	"github.com/kougami132/MsgPilot/internal/channels"
	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
)

type HandlerService interface {
	CommonPush(ticket string, title string, msg string) (*models.Message, error)
	OneBotPush(ticket string, msg string) (*models.Message, error)
}

type handlerService struct {
	bridgeService  BridgeService
	messageService MessageService
}

func NewHandlerService(bridgeService BridgeService, messageService MessageService) HandlerService {
	return &handlerService{bridgeService: bridgeService, messageService: messageService}
}

// processPush 是一个处理通用消息发送逻辑的私有辅助函数
func (u *handlerService) processPush(
	ticket string,
	expectedSourceType types.ChannelType,
	createMessageFunc func(bridge *models.Bridge) *models.Message,
) (*models.Message, error) {
	bridge, err := u.bridgeService.GetBridgeByTicket(ticket)
	if err != nil {
		return nil, err
	}

	if bridge.SourceChannelType != expectedSourceType {
		return nil, errors.New("中转源渠道不匹配")
	}

	if !bridge.IsActive {
		return nil, errors.New("中转配置未激活")
	}

	message := createMessageFunc(bridge)

	err = u.messageService.CreateMessage(message)
	if err != nil {
		return nil, err
	}

	// 发送消息
	go func() {
		u.messageService.UpdateMessageStatus(message, types.StatusSending)
		handler, err := channels.GetChannelHandler(bridge.TargetChannel)
		if err != nil {
			u.messageService.UpdateMessageStatusWithErrorMessage(message, types.StatusFailed, "中转目标渠道不可用")
			return
		}
		err = handler.Send(message)
		if err != nil {
			u.messageService.UpdateMessageStatusWithErrorMessage(message, types.StatusFailed, err.Error())
			return
		}
		u.messageService.UpdateMessageStatus(message, types.StatusSuccess)
	}()

	return message, nil
}

func (u *handlerService) OneBotPush(ticket string, msg string) (*models.Message, error) {
	createFunc := func(bridge *models.Bridge) *models.Message {
		return &models.Message{
			Title:    "MsgPilot消息推送",
			Content:  msg,
			Status:   types.StatusPending,
			BridgeID: bridge.ID,
			Bridge:   *bridge,
		}
	}
	return u.processPush(ticket, types.TypeOneBot, createFunc)
}

func (u *handlerService) CommonPush(ticket string, title string, body string) (*models.Message, error) {
	createFunc := func(bridge *models.Bridge) *models.Message {
		return &models.Message{
			Title:    title,
			Content:  body,
			Status:   types.StatusPending,
			BridgeID: bridge.ID,
			Bridge:   *bridge,
		}
	}
	return u.processPush(ticket, types.TypeBark, createFunc)
}