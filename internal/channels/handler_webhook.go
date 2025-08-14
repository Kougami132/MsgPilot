package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/internal/utils"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

type WebhookHandler struct {
	config datatypes.JSON
}

type WebhookConfig struct {
	BaseUrl     string            `json:"base_url"`
	Method      string            `json:"method"`
	ContentType string            `json:"content_type"`
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
}

func (h *WebhookHandler) Send(message *models.Message) error {
	var cfg WebhookConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析Webhook配置失败: %w", err)
	}

	body := strings.Replace(cfg.Body, "$title", message.Title, -1)
	body = strings.Replace(body, "$content", message.Content, -1)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}

	req, err := http.NewRequest(cfg.Method, cfg.BaseUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", cfg.ContentType)
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelAdapter(types.TypeWebhook)
	RegisterChannelHandler(types.TypeWebhook, func(config datatypes.JSON) ChannelHandler {
		return &WebhookHandler{config: config}
	})
}
