package usecase

import (
	"errors"

	"github.com/kougami132/MsgPilot/channels"
	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
)

type AdapterUsecase interface {
	OneBotSendMessage(ticket string, msg string) (*models.Message, error)
}

type adapterUsecase struct {
	bridgeUsecase  BridgeUsecase
	messageUsecase MessageUsecase
}

func NewAdapterUsecase(bridgeUsecase BridgeUsecase, messageUsecase MessageUsecase) AdapterUsecase {
	return &adapterUsecase{bridgeUsecase: bridgeUsecase, messageUsecase: messageUsecase}
}

func (u *adapterUsecase) OneBotSendMessage(ticket string, msg string) (*models.Message, error) {
	bridge, err := u.bridgeUsecase.GetBridgeByTicket(ticket)
	if err != nil {
		return nil, err
	}

	if bridge.SourceChannelType != types.TypeOneBot {
		return nil, errors.New("中转源渠道不匹配")
	}

	if !bridge.IsActive {
		return nil, errors.New("中转配置未激活")
	}

	// 转换成通用消息
	message := &models.Message{
		Title:    "MsgPilot消息推送",
		Content:  msg,
		Status:   types.StatusPending,
		BridgeID: bridge.ID,
		Bridge:   *bridge,
	}
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
