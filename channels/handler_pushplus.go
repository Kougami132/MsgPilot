package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

type PushPlusHandler struct {
	config datatypes.JSON
}

type PushPlusConfig struct {
	Token string `json:"token"`
}

func (h *PushPlusHandler) Send(message *models.Message) error {
	var cfg PushPlusConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析PushPlus配置失败: %w", err)
	}

	body := map[string]interface{}{
		"token": cfg.Token,
		"title": message.Title,
		"content": message.Content,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post("http://www.pushplus.plus/send", "application/json;charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func init() {
	RegisterChannelHandler(types.TypePushPlus, func(config datatypes.JSON) ChannelHandler {
		return &PushPlusHandler{config: config}
	})
}

