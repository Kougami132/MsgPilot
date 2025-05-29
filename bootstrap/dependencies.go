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
	ConfigController  *controller.ConfigController
	MessageController *controller.MessageController
	BridgeController  *controller.BridgeController
}

// NewAppDependencies 初始化应用的所有依赖项
func NewAppDependencies(app Application) *AppDependencies {
	// 初始化仓库
	userRepo := repository.NewUserRepository(app.SQLite.DB)
	channelRepo := repository.NewChannelRepository(app.SQLite.DB)
	configRepo := repository.NewConfigRepository(app.SQLite.DB)
	messageRepo := repository.NewMessageRepository(app.SQLite.DB)
	bridgeRepo := repository.NewBridgeRepository(app.SQLite.DB)

	// 初始化用例
	authUseCase := usecase.NewAuthUseCase(userRepo, app.Env)
	channelUseCase := usecase.NewChannelUsecase(channelRepo)
	configUseCase := usecase.NewConfigUsecase(configRepo)
	messageUseCase := usecase.NewMessageUsecase(messageRepo)
	bridgeUsecase := usecase.NewBridgeUsecase(bridgeRepo, channelRepo)

	// 初始化config数据
	// configUseCase.InitConfig("username", "admin")

	// 初始化控制器
	authController := controller.NewAuthController(authUseCase, app.Env)
	channelController := controller.NewChannelController(channelUseCase)
	configController := controller.NewConfigController(configUseCase)
	messageController := controller.NewMessageController(messageUseCase)
	bridgeController := controller.NewBridgeController(bridgeUsecase)

	return &AppDependencies{
		AuthController:    authController,
		ChannelController: channelController,
		ConfigController:  configController,
		MessageController: messageController,
		BridgeController:  bridgeController,
	}
}
