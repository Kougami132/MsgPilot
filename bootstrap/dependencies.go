package bootstrap

import (
	"github.com/kougami132/MsgPilot/api/controller"
	"github.com/kougami132/MsgPilot/internal/repository"
	"github.com/kougami132/MsgPilot/internal/service"
)

// AppDependencies 包含应用的所有依赖项
type AppDependencies struct {
	AuthController    *controller.AuthController
	ChannelController *controller.ChannelController
	MessageController *controller.MessageController
	BridgeController  *controller.BridgeController
	AdapterController *controller.AdapterController
}

// NewAppDependencies 初始化应用的所有依赖项
func NewAppDependencies(app Application) *AppDependencies {
	// 初始化仓库
	userRepo := repository.NewUserRepository(app.SQLite.DB)
	channelRepo := repository.NewChannelRepository(app.SQLite.DB)
	messageRepo := repository.NewMessageRepository(app.SQLite.DB)
	bridgeRepo := repository.NewBridgeRepository(app.SQLite.DB)

	// 初始化用例
	authService := service.NewAuthService(userRepo, app.Env)
	channelService := service.NewChannelService(channelRepo)
	messageService := service.NewMessageService(messageRepo)
	bridgeService := service.NewBridgeService(bridgeRepo, channelRepo)
	handlerService := service.NewHandlerService(bridgeService, messageService)

	// 初始化控制器
	authController := controller.NewAuthController(authService, app.Env)
	channelController := controller.NewChannelController(channelService)
	messageController := controller.NewMessageController(messageService)
	bridgeController := controller.NewBridgeController(bridgeService)
	adapterController := controller.NewAdapterController(handlerService)

	return &AppDependencies{
		AuthController:    authController,
		ChannelController: channelController,
		MessageController: messageController,
		BridgeController:  bridgeController,
		AdapterController: adapterController,
	}
}
