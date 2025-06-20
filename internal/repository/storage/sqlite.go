package storage

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite" // 纯Go实现的SQLite驱动
	"gorm.io/gorm"
)

type SQLiteDB struct {
	DB *gorm.DB
}

// NewSQLiteDB 创建一个新的SQLite数据库连接
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	// 确保数据库目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 确保关闭外键约束
	})
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return &SQLiteDB{DB: db}, nil
}

// AutoMigrate 自动迁移数据库模型
func (db *SQLiteDB) AutoMigrate(models ...interface{}) error {
	return db.DB.AutoMigrate(models...)
}

// Close 关闭SQLite数据库连接
func (db *SQLiteDB) Close() error {
	if db.DB != nil {
		sqlDB, err := db.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Transaction 执行事务
func (db *SQLiteDB) Transaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}
