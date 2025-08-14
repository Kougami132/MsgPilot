package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/internal/utils"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

type TelegramHandler struct {
	config datatypes.JSON
}

type TelegramConfig struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

func (h *TelegramHandler) Send(message *models.Message) error {
	var cfg TelegramConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析Telegram配置失败: %w", err)
	}

	body := map[string]interface{}{
		"chat_id":                  cfg.UserId,
		"text":                     fmt.Sprintf("%s：\n%s", message.Title, message.Content),
		"disable_web_page_preview": true,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.Token), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelHandler(types.TypeTelegram, func(config datatypes.JSON) ChannelHandler {
		return &TelegramHandler{config: config}
	})
}
