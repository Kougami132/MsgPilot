package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/datatypes"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/internal/utils"
	"github.com/kougami132/MsgPilot/models"
)

type WxPusherHandler struct {
	config datatypes.JSON
}

type WxPusherConfig struct {
	Token string `json:"token"`
	Topic string `json:"topic"`
	Uid   string `json:"uid"`
}

func (h *WxPusherHandler) Send(message *models.Message) error {
	var cfg WxPusherConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析WxPusher配置失败: %w", err)
	}

	body := map[string]interface{}{
		"contentType": 1,
		"summary":     message.Title,
		"content":     message.Content,
		"topicIds":    strings.Split(cfg.Topic, ","),
		"uids":        strings.Split(cfg.Uid, ","),
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}

	resp, err := http.Post("https://wxpusher.zjiecode.com/api/send/message", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelHandler(types.TypeWxPusher, func(config datatypes.JSON) ChannelHandler {
		return &WxPusherHandler{config: config}
	})
}
