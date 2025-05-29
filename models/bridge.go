package models

import (
	"time"
)

// Bridge 桥接模型
type Bridge struct {
	ID              string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	SourceChannelID string    `gorm:"type:char(36);not null"`
	SourceChannel   Channel   `gorm:"foreignKey:SourceChannelID"`
	TargetChannelID string    `gorm:"type:char(36);not null"`
	TargetChannel   Channel   `gorm:"foreignKey:TargetChannelID"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Bridge) TableName() string {
	return "bridges"
}
