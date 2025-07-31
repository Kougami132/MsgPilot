package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/api/middleware"
	"github.com/kougami132/MsgPilot/bootstrap"

	_ "github.com/kougami132/MsgPilot/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// 静态文件服务 - 提供Vue构建后的文件
	router.Static("/assets", "./frontend/dist/assets")
	router.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")

	// 首页路由
	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// API路由组
	apiGroup := router.Group("/api")

	// 添加 Swagger 路由
	apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 适配器路由组
	deps.AdapterController.RegisterRoutes(apiGroup)

	// 认证相关路由 (不需要 Token 验证)
	deps.AuthController.RegisterRoutes(apiGroup) // AuthController 的路由注册在 /api/auth 下

	// 以下路由组需要 Token 认证
	protectedRoutes := apiGroup.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware(app))

	// 为其他 Controller 注册路由到 protectedRoutes
	deps.ChannelController.RegisterRoutes(protectedRoutes)
	deps.BridgeController.RegisterRoutes(protectedRoutes)
	deps.MessageController.RegisterRoutes(protectedRoutes)

	// 处理Vue Router的history模式，所有非API路由都返回index.html
	router.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		// 否则返回Vue应用的index.html（支持前端路由）
		c.File("./frontend/dist/index.html")
	})

	return router
}
