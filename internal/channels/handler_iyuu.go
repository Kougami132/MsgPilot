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

type IYUUHandler struct {
	config datatypes.JSON
}

type IYUUConfig struct {
	Token string `json:"token"`
}

func (h *IYUUHandler) Send(message *models.Message) error {
	var cfg IYUUConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析IYUU配置失败: %w", err)
	}

	body := map[string]interface{}{
		"text": message.Title,
		"desp": message.Content,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("https://iyuu.co/%s.send", cfg.Token), "application/json;charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelHandler(types.TypeIYUU, func(config datatypes.JSON) ChannelHandler {
		return &IYUUHandler{config: config}
	})
}
