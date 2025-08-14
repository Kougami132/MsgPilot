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

type QMsgHandler struct {
	config datatypes.JSON
}

type QMsgConfig struct {
	Token string `json:"token"`
	Type  string `json:"type"`
	Bot   string `json:"bot"`
	QQ    string `json:"qq"`
}

func (h *QMsgHandler) Send(message *models.Message) error {
	var cfg QMsgConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析QMsg配置失败: %w", err)
	}

	body := map[string]interface{}{
		"bot": cfg.Bot,
		"qq":  cfg.QQ,
		"msg": fmt.Sprintf("%s：\n%s", message.Title, message.Content),
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}
	resp, err := http.Post(fmt.Sprintf("https://qmsg.zendee.cn/j%s/%s", cfg.Type, cfg.Token), "application/json;charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelHandler(types.TypeQMSG, func(config datatypes.JSON) ChannelHandler {
		return &QMsgHandler{config: config}
	})
}
