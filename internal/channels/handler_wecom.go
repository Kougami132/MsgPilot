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

type WeComHandler struct {
	config datatypes.JSON
}

type WeComConfig struct {
	Token   string `json:"token"`
	CorpId  string `json:"corp_id"`
	AgentId string `json:"agent_id"`
	Secret  string `json:"secret"`
}

func (h *WeComHandler) Send(message *models.Message) error {
	var cfg WeComConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析WeCom配置失败: %w", err)
	}

	resp, err := http.Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", cfg.CorpId, cfg.Secret))
	if err != nil {
		return fmt.Errorf("获取access_token失败: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析access_token响应失败: %w", err)
	}

	accessToken := result["access_token"].(string)

	body := map[string]interface{}{
		"touser":  "@all",
		"agentid": cfg.AgentId,
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": fmt.Sprintf("%s：\n%s", message.Title, message.Content),
		},
		"duplicate_check_interval": 600,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}

	resp, err = http.Post(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", accessToken), "application/json;charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return utils.CheckHTTPResponse(resp)
}

func init() {
	RegisterChannelHandler(types.TypeWeCom, func(config datatypes.JSON) ChannelHandler {
		return &WeComHandler{config: config}
	})
}
