package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/api/middleware"
	"github.com/kougami132/MsgPilot/bootstrap"
)

// Setup 配置和返回 Gin 路由器
func Setup(app bootstrap.Application, deps *bootstrap.AppDependencies) *gin.Engine {
	// 根据环境设置Gin模式
	if app.Env.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由
	router := gin.Default()

	// 注册中间件
	router.Use(gin.Logger(), gin.Recovery())

	// 注册根路由
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "MsgPilot API 服务正在运行",
		})
	})

	// API路由组
	apiGroup := router.Group("/api")

	// 认证相关路由 (不需要 Token 验证)
	deps.AuthController.RegisterRoutes(apiGroup) // AuthController 的路由注册在 /api/auth 下

	// 以下路由组需要 Token 认证
	protectedRoutes := apiGroup.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware(app))

	// 为其他 Controller 注册路由到 protectedRoutes  
	deps.ChannelController.RegisterRoutes(protectedRoutes) 
	deps.ConfigController.RegisterRoutes(protectedRoutes)  
	deps.MessageController.RegisterRoutes(protectedRoutes) 

	return router
}
