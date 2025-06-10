package usecase

import (
	"errors"

	"github.com/kougami132/MsgPilot/channels"
	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
)

type HandlerUsecase interface {
	CommonPush(ticket string, title string, msg string) (*models.Message, error)
	OneBotPush(ticket string, msg string) (*models.Message, error)
}

type handlerUsecase struct {
	bridgeUsecase  BridgeUsecase
	messageUsecase MessageUsecase
}

func NewHandlerUsecase(bridgeUsecase BridgeUsecase, messageUsecase MessageUsecase) HandlerUsecase {
	return &handlerUsecase{bridgeUsecase: bridgeUsecase, messageUsecase: messageUsecase}
}

// processPush 是一个处理通用消息发送逻辑的私有辅助函数
func (u *handlerUsecase) processPush(
	ticket string,
	expectedSourceType types.ChannelType,
	createMessageFunc func(bridge *models.Bridge) *models.Message,
) (*models.Message, error) {
	bridge, err := u.bridgeUsecase.GetBridgeByTicket(ticket)
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

	err = u.messageUsecase.CreateMessage(message)
	if err != nil {
		return nil, err
	}

	// 发送消息
	go func() {
		u.messageUsecase.UpdateMessageStatus(message, types.StatusSending)
		handler, err := channels.GetChannelHandler(bridge.TargetChannel)
		if err != nil {
			u.messageUsecase.UpdateMessageStatusWithErrorMessage(message, types.StatusFailed, "中转目标渠道不可用")
			return
		}
		err = handler.Send(message)
		if err != nil {
			u.messageUsecase.UpdateMessageStatusWithErrorMessage(message, types.StatusFailed, err.Error())
			return
		}
		u.messageUsecase.UpdateMessageStatus(message, types.StatusSuccess)
	}()

	return message, nil
}

func (u *handlerUsecase) OneBotPush(ticket string, msg string) (*models.Message, error) {
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

func (u *handlerUsecase) CommonPush(ticket string, title string, body string) (*models.Message, error) {
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