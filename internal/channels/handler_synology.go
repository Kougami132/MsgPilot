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

type SynologyHandler struct {
	config datatypes.JSON
}

type SynologyConfig struct {
	BaseUrl string `json:"base_url"`
}

func (h *SynologyHandler) Send(message *models.Message) error {
	var cfg SynologyConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析Synology配置失败: %w", err)
	}

	body := map[string]interface{}{
		"text": fmt.Sprintf("%s：\n%s", message.Title, message.Content),
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(cfg.BaseUrl, "application/json;charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelHandler(types.TypeSynology, func(config datatypes.JSON) ChannelHandler {
		return &SynologyHandler{config: config}
	})
}
