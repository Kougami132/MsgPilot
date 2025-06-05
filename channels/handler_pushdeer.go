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

type PushDeerHandler struct {
	config datatypes.JSON
}

type PushDeerConfig struct {
	BaseUrl string `json:"base_url"`
	PushKey string `json:"pushkey"`
}

func (h *PushDeerHandler) Send(message *models.Message) error {
	var cfg PushDeerConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析PushDeer配置失败: %w", err)
	}

	body := map[string]interface{}{
		"pushkey": cfg.PushKey,
		"text": message.Title,
		"desp": message.Content,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/message/push", cfg.BaseUrl), "application/json;charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func init() {
	RegisterChannelAdapter(types.TypePushDeer)
	RegisterChannelHandler(types.TypePushDeer, func(config datatypes.JSON) ChannelHandler {
		return &PushDeerHandler{config: config}
	})
}

