package models

import (
	"time"
)

// 配置模型
type Config struct {
	Key       string    `gorm:"type:varchar(255);primaryKey"` // 配置键
	Value     string    `gorm:"type:text"`                   // 配置值
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定表名
func (Config) TableName() string {
	return "configs"
} 