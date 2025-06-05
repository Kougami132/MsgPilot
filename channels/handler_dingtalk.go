package channels

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kougami132/MsgPilot/internal/types"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

type DingTalkHandler struct {
	config datatypes.JSON
}

type DingTalkConfig struct {
	Token string `json:"token"`
	Secret string `json:"secret"`
}

func (h *DingTalkHandler) Send(message *models.Message) error {
	var cfg DingTalkConfig
	if err := json.Unmarshal(h.config, &cfg); err != nil {
		return fmt.Errorf("解析钉钉机器人配置失败: %w", err)
	}

	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("%s：\n%s", message.Title, message.Content),
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息体失败: %w", err)
	}

	var base_url string
	if cfg.Secret == "" {
		base_url = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", cfg.Token)
	} else {
		timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
		stringToSign := fmt.Sprintf("%s\n%s", timestamp, cfg.Secret)
		h := hmac.New(sha256.New, []byte(cfg.Secret))
		h.Write([]byte(stringToSign))
		sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
		base_url = fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%s&sign=%s",
			cfg.Token, timestamp, url.QueryEscape(sign))
	}
	resp, err := http.Post(base_url, "application/json; charset=utf-8", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func init() {
	RegisterChannelHandler(types.TypeDingTalk, func(config datatypes.JSON) ChannelHandler {
		return &DingTalkHandler{config: config}
	})
}