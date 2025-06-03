package bootstrap

import (
	"log"
	"path/filepath"

	"github.com/kougami132/MsgPilot/config"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/sqlite"
)

type Application struct {
	Env    *config.Env
	SQLite *sqlite.SQLiteDB
}

func App() Application {
	app := &Application{}
	app.Env = config.NewEnv()

	// 初始化SQLite数据库
	dbPath := filepath.Join("data", "msgpilot.db")
	sqliteDB, err := sqlite.NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatalf("SQLite初始化失败: %v", err)
	}

	// 自动迁移数据库模型
	if err := sqliteDB.AutoMigrate(&models.User{}, &models.Message{}, &models.Bridge{}, &models.Channel{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移成功")

	app.SQLite = sqliteDB

	return *app
}

// Close 关闭应用资源
func (app *Application) Close() {
	if app.SQLite != nil {
		if err := app.SQLite.Close(); err != nil {
			log.Printf("关闭SQLite数据库失败: %v", err)
		}
	}
}
