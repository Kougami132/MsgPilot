package models

import (
	"time"

	"gorm.io/gorm"
)

// 状态枚举定义
type MessageStatus int
const (
    StatusPending MessageStatus = iota // 待发送
    StatusSending                      // 发送中
    StatusSent                         // 已发送
    StatusFailed                       // 失败
)

// Message 消息模型
type Message struct {
	ID          	 string         `gorm:"type:char(36);primaryKey"` // UUID
	Content     	 string         `gorm:"type:text;not null"`       // 消息内容
	Status      	 MessageStatus  `gorm:"index;default:0"`          // 状态枚举
	SourceChannelID  string         `gorm:"type:char(36);not null"`  // 接收渠道ID
	SourceChannel    Channel        `gorm:"foreignKey:SourceChannelID"` // 接收渠道
	TargetChannelID  string         `gorm:"type:char(36);not null"`  // 发送渠道ID
	TargetChannel    Channel        `gorm:"foreignKey:TargetChannelID"` // 发送渠道
	CreatedAt 		 time.Time      `json:"created_at"`
	UpdatedAt 		 time.Time      `json:"updated_at"`
	DeletedAt 		 gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
