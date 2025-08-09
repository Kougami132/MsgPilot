package models

import (
	"time"

	"github.com/kougami132/MsgPilot/internal/types"
	"gorm.io/datatypes"
)

// Channel 渠道模型
type Channel struct {
	ID        int               `gorm:"primaryKey"                      json:"id"`                          // UUID
	Name      string            `gorm:"type:varchar(100);uniqueIndex"   json:"name"`                        // 渠道名称（如：Email-SMTP、OneBot-v11）
	Type      types.ChannelType `gorm:"index"                           json:"type"`                        // 渠道类型
	Config    datatypes.JSON    `gorm:"type:json"                       json:"config" swaggertype:"object"` // 渠道配置（统一结构）
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// TableName 指定表名
func (Channel) TableName() string {
	return "channels"
}
