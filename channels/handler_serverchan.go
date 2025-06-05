package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
	"github.com/kougami132/MsgPilot/internal/types"
)

type ServerChanHandler struct {
	config datatypes.JSON
}

type ServerChanConfig struct {
	SendKey string `json:"sendkey"`
}

func (h *ServerChanHandler) Send(message *models.Message) error {
	var cfg ServerChanConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析ServerChan配置失败: %w", err)
	}

	body := map[string]interface{}{
		"title": message.Title,
		"desp": message.Content,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("https://sctapi.ftqq.com/%s.send", cfg.SendKey), "application/json;charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func init() {
	RegisterChannelHandler(types.TypeServerChan, func(config datatypes.JSON) ChannelHandler {
		return &ServerChanHandler{config: config}
	})
}
