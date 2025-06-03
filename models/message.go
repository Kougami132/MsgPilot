package models

import (
	"time"

	"github.com/kougami132/MsgPilot/internal/types"
	"gorm.io/gorm"
)

// Message 消息模型
type Message struct {
	ID          	 	int         		`gorm:"primaryKey"                  json:"id"`
	Title     	 		string          	`gorm:"type:text;not null"          json:"title"`
	Content     	 	string          	`gorm:"type:text;not null"          json:"content"`
	Status      	 	types.MessageStatus `gorm:"index;default:0"             json:"status"`
	ErrorMessage 		string          	`gorm:"type:text"                   json:"error_message"`
	BridgeID 		 	int         		`gorm:"not null"                    json:"bridge_id"`
	Bridge 			 	Bridge         		`gorm:"foreignKey:BridgeID"         json:"bridge"`
	CreatedAt 		 	time.Time       	`json:"created_at"`
	UpdatedAt 		 	time.Time       	`json:"updated_at"`
	DeletedAt 		 	gorm.DeletedAt  	`gorm:"index" json:"-"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
