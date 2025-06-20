package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/datatypes"

	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/internal/types"
)

type GotifyHandler struct {
	config datatypes.JSON
}

type GotifyConfig struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

func (h *GotifyHandler) Send(message *models.Message) error {
	var cfg GotifyConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析Gotify配置失败: %w", err)
	}
	
	body := map[string]interface{}{
		"title":   	message.Title,
		"message":  message.Content,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/message?token=%s", cfg.BaseURL, cfg.Token), "application/json; charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func init() {
	RegisterChannelAdapter(types.TypeGotify)
	RegisterChannelHandler(types.TypeGotify, func(config datatypes.JSON) ChannelHandler {
		return &GotifyHandler{config: config}
	})
}

