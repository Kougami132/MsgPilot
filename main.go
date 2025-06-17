package main

import (
	"github.com/kougami132/MsgPilot/bootstrap"
	"github.com/kougami132/MsgPilot/route"
)

func main() {
	// 1. 初始化应用核心
	app := bootstrap.App()
	defer app.Close()

	// 2. 初始化所有依赖项 (仓库、用例、控制器)
	dependencies := bootstrap.NewAppDependencies(app)

	// 3. 设置路由
	router := route.Setup(app, dependencies)

	// 4. 启动服务器
	router.Run(":" + app.Env.Port)
}
