package bootstrap

import (
	"github.com/kougami132/MsgPilot/api/controller"
	"github.com/kougami132/MsgPilot/repository"
	"github.com/kougami132/MsgPilot/usecase"
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
	authUseCase := usecase.NewAuthUseCase(userRepo, app.Env)
	channelUseCase := usecase.NewChannelUsecase(channelRepo)
	messageUseCase := usecase.NewMessageUsecase(messageRepo)
	bridgeUsecase := usecase.NewBridgeUsecase(bridgeRepo, channelRepo)
	handlerUsecase := usecase.NewHandlerUsecase(bridgeUsecase, messageUseCase)

	// 初始化控制器
	authController := controller.NewAuthController(authUseCase, app.Env)
	channelController := controller.NewChannelController(channelUseCase)
	messageController := controller.NewMessageController(messageUseCase)
	bridgeController := controller.NewBridgeController(bridgeUsecase)
	adapterController := controller.NewAdapterController(handlerUsecase)

	return &AppDependencies{
		AuthController:    authController,
		ChannelController: channelController,
		MessageController: messageController,
		BridgeController:  bridgeController,
		AdapterController: adapterController,
	}
}
