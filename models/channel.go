package models

import (
	"time"

	"gorm.io/datatypes"
	"github.com/kougami132/MsgPilot/internal/types"
)

// 渠道模型
type Channel struct {
    ID          string              `gorm:"type:char(36);primaryKey"` // UUID
    Name        string              `gorm:"type:varchar(100);uniqueIndex"` // 渠道名称（如：Email-SMTP、OneBot-v11）
    Type        types.ChannelType   `gorm:"index"`      // 渠道类型
    Direction   types.DirectionType `gorm:"index"`      // 方向类型
    Config      datatypes.JSON      `gorm:"type:json"`  // 渠道配置（统一结构）
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// TableName 指定表名
func (Channel) TableName() string {
	return "channels"
}