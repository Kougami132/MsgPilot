package channels

import (
	"fmt"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

var channelRegistry = make(map[types.ChannelType]func(config datatypes.JSON) ChannelHandler)

// 注册渠道类型
func RegisterChannelHandler(
    typeName types.ChannelType, 
    creator func(config datatypes.JSON) ChannelHandler,
) {
    channelRegistry[typeName] = creator
}

// 获取渠道类型
func GetChannelHandler(channel models.Channel) (ChannelHandler, error) {
    creator, ok := channelRegistry[channel.Type]
    if !ok {
        return nil, fmt.Errorf("未注册的渠道类型: %s", channel.Type)
    }
    return creator(channel.Config), nil
}
