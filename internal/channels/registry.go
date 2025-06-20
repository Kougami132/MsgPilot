package channels

import (
	"fmt"
	"sort"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

var channelAdapterRegistry = make([]types.ChannelType, 0)
var channelHandlerRegistry = make(map[types.ChannelType]func(config datatypes.JSON) ChannelHandler)

// 注册渠道适配器
func RegisterChannelAdapter(adapter types.ChannelType) {
	channelAdapterRegistry = append(channelAdapterRegistry, adapter)
}

// 注册渠道处理器
func RegisterChannelHandler(
	typeName types.ChannelType,
	creator func(config datatypes.JSON) ChannelHandler,
) {
	channelHandlerRegistry[typeName] = creator
}

// 获取所有已注册的渠道适配器
func GetChannelAdapters() []types.ChannelType {
	return channelAdapterRegistry
}

// 获取所有已注册的渠道处理器
func GetChannelHandlers() []types.ChannelType {
	types := make([]types.ChannelType, 0, len(channelHandlerRegistry))
	for typeName := range channelHandlerRegistry {
		types = append(types, typeName)
	}
	// 排序以固定每次返回的顺序
	sort.Slice(types, func(i, j int) bool {
		return types[i] < types[j]
	})
	return types
}

// 取出渠道处理器
func GetChannelHandler(channel models.Channel) (ChannelHandler, error) {
	creator, ok := channelHandlerRegistry[channel.Type]
	if !ok {
		return nil, fmt.Errorf("未注册的渠道类型: %s", channel.Type)
	}
	return creator(channel.Config), nil
}
