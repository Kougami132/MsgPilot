package models

import (
	"time"

	"github.com/kougami132/MsgPilot/internal/types"
)

// Bridge 桥接模型
type Bridge struct {
	ID                	int       			`gorm:"primaryKey"                  json:"id"`
	Name              	string    			`gorm:"size:100;not null"           json:"name"`
	Ticket            	string    			`gorm:"size:100;not null"           json:"ticket"`
	SourceChannelType 	types.ChannelType   `gorm:"index"                       json:"source_channel_type"`
	TargetChannelID   	int       			`gorm:"not null"                    json:"target_channel_id"`
	TargetChannel     	Channel   			`gorm:"foreignKey:TargetChannelID"  json:"target_channel"`
	IsActive          	bool      			`gorm:"default:true"                json:"is_active"`
	CreatedAt         	time.Time 			`json:"created_at"`
	UpdatedAt         	time.Time 			`json:"updated_at"`
}

func (Bridge) TableName() string {
	return "bridges"
}
