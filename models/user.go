package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int            `gorm:"primarykey"               json:"id"`
	Username  string         `gorm:"size:100;not null;unique" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

