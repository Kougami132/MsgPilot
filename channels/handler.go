package channels

import "github.com/kougami132/MsgPilot/models"

type ChannelHandler interface {
    Send(message *models.Message) error
}