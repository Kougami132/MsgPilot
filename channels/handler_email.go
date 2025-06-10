package channels

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

type EmailHandler struct {
	config datatypes.JSON
}

type EmailConfig struct {
	Host     	string `json:"host"`
	Port     	int    `json:"port"`
	FromEmail   string `json:"from_email"`
	ToEmail    	string `json:"to_email"`
	AuthCode 	string `json:"auth_code"`
}

func (h *EmailHandler) Send(message *models.Message) error {
	var cfg EmailConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析Email配置失败: %w", err)
	}

	em := email.NewEmail()
	em.From = cfg.FromEmail
	em.To = []string{cfg.ToEmail}
	em.Subject = message.Title
	em.Text = []byte(message.Content)
	err := em.Send(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), smtp.PlainAuth("", cfg.FromEmail, cfg.AuthCode, cfg.Host))
	if err != nil {
		return fmt.Errorf("发送Email失败: %w", err)
	}
	return nil
}

func init() {
	RegisterChannelHandler(types.TypeEmail, func(config datatypes.JSON) ChannelHandler {
		return &EmailHandler{config: config}
	})
}