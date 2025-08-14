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

type FeishuHandler struct {
	config datatypes.JSON
}

type FeishuConfig struct {
	FSKey string `json:"fs_key"`
}

func (h *FeishuHandler) Send(message *models.Message) error {
	var cfg FeishuConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析飞书机器人配置失败: %w", err)
	}

	body := map[string]interface{}{
		"msg_type": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("%s：\n%s", message.Title, message.Content),
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", cfg.FSKey), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelHandler(types.TypeFeiShu, func(config datatypes.JSON) ChannelHandler {
		return &FeishuHandler{config: config}
	})
}
