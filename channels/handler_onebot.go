package channels

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/datatypes"
)

type OneBotHandler struct {
	config datatypes.JSON
}

// OneBot 配置参数。
type OneBotConfig struct {
	BaseUrl     string `json:"base_url"` // BaseUrl 是 OneBot API 的基础 URL，例如 http://127.0.0.1:5700/send_msg。
	MessageType string `json:"message_type"` // MessageType 指定消息类型，可以是 "private" (私聊) 或 "group" (群聊)。
	UserID      string `json:"user_id,omitempty"` // UserID 是私聊消息的目标用户 ID。仅当 MessageType 为 "private" 时需要。
	GroupID     string `json:"group_id,omitempty"` // GroupID 是群聊消息的目标群组 ID。仅当 MessageType 为 "group" 时需要。
}

func (h *OneBotHandler) Send(content string) error {
	// 解析存储在 h.config 中的 JSON。
	var cfg OneBotConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析OneBot配置失败: %w", err)
	}

	switch cfg.MessageType {
	case "private":
		if cfg.UserID == "" {
			return fmt.Errorf("private 消息类型缺少 user_id")
		}
		// 构造向 OneBot API 发送私聊消息的 URL 并发送 GET 请求。
		resp, err := http.Get(cfg.BaseUrl + "?message_type=private&user_id=" + cfg.UserID + "&message=" + content)
		if err != nil {
			return fmt.Errorf("发送OneBot私聊消息失败: %w", err)
		}
		defer resp.Body.Close()
	case "group":
		if cfg.GroupID == "" {
			return fmt.Errorf("group 消息类型缺少 group_id")
		}
		// 构造向 OneBot API 发送群聊消息的 URL 并发送 GET 请求。
		resp, err := http.Get(cfg.BaseUrl + "?message_type=group&group_id=" + cfg.GroupID + "&message=" + content)
		if err != nil {
			return fmt.Errorf("发送OneBot群聊消息失败: %w", err)
		}
		defer resp.Body.Close()
	default:
		return fmt.Errorf("未知的OneBot message_type: %s", cfg.MessageType)
	}

	return nil
}

func init() {
	RegisterChannelType("onebot", func(config datatypes.JSON) ChannelHandler {
		return &OneBotHandler{config: config}
	})
}
