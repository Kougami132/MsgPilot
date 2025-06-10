package channels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"gorm.io/datatypes"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/internal/types"
)

type NtfyHandler struct {
	config datatypes.JSON
}

type NtfyConfig struct {
	BaseUrl string `json:"base_url"`
	Token   string `json:"token"`
	Topic   string `json:"topic"`
}

func (h *NtfyHandler) Send(message *models.Message) error {
	var cfg NtfyConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析Ntfy配置失败: %w", err)
	}
	
	body := map[string]interface{}{
		"topic":   	cfg.Topic,
		"title":   	message.Title,
		"message":  message.Content,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化Ntfy消息失败: %w", err)
	}

	resp, err := http.Post(cfg.BaseUrl, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("发送Ntfy请求失败: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func init() {
	RegisterChannelAdapter(types.TypeNtfy)
	RegisterChannelHandler(types.TypeNtfy, func(config datatypes.JSON) ChannelHandler {
		return &NtfyHandler{config: config}
	})
}

