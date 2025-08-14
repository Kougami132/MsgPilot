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

type OneBotHandler struct {
	config datatypes.JSON
}

// OneBot 配置参数。
type OneBotConfig struct {
	BaseUrl     string `json:"base_url"`           // BaseUrl 是 OneBot API 的基础 URL，例如 http://127.0.0.1:5700/send_msg。
	MessageType string `json:"message_type"`       // MessageType 指定消息类型，可以是 "private" (私聊) 或 "group" (群聊)。
	UserID      string `json:"user_id,omitempty"`  // UserID 是私聊消息的目标用户 ID。仅当 MessageType 为 "private" 时需要。
	GroupID     string `json:"group_id,omitempty"` // GroupID 是群聊消息的目标群组 ID。仅当 MessageType 为 "group" 时需要。
}

func (h *OneBotHandler) Send(message *models.Message) error {
	// 解析存储在 h.config 中的 JSON。
	var cfg OneBotConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析OneBot配置失败: %w", err)
	}

	body := map[string]interface{}{
		"message_type": cfg.MessageType,
		"user_id":      cfg.UserID,
		"group_id":     cfg.GroupID,
		"message":      message.Title + "：\n" + message.Content,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(cfg.BaseUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelAdapter(types.TypeOneBot)
	RegisterChannelHandler(types.TypeOneBot, func(config datatypes.JSON) ChannelHandler {
		return &OneBotHandler{config: config}
	})
}
