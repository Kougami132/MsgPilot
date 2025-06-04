package channels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

type BarkHandler struct {
	config datatypes.JSON
}

type BarkConfig struct {
	BaseUrl string `json:"base_url"`
	Key     string `json:"key"`
	Sound   string `json:"sound"`
	Icon    string `json:"icon"`
	Url     string `json:"url"`
}

func (h *BarkHandler) Send(message *models.Message) error {
	var cfg BarkConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析Bark配置失败: %w", err)
	}

	body := map[string]interface{}{
		"title":   message.Title,
		"body":    message.Content,
		"sound":   cfg.Sound,
		"icon":    cfg.Icon,
		"url":     cfg.Url,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(cfg.BaseUrl, "application/json; charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func init() {
	RegisterChannelHandler(types.TypeBark, func(config datatypes.JSON) ChannelHandler {
		return &BarkHandler{config: config}
	})
}
