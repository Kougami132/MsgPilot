package types

// 渠道类型
type ChannelType string
const (
	TypeBark       ChannelType = "Bark"       // Bark
	TypeDingTalk   ChannelType = "DingTalk"   // 钉钉机器人
	TypeFeiShu     ChannelType = "FeiShu"     // 飞书机器人
	TypeOneBot     ChannelType = "OneBot"     // OneBot
	TypeGotify     ChannelType = "Gotify"     // Gotify
	TypeServerChan ChannelType = "ServerChan" // Server酱
	TypePushDeer   ChannelType = "PushDeer"   // PushDeer
	TypeSynology   ChannelType = "Synology"   // 群晖chat
	TypePushPlus   ChannelType = "PushPlus"   // PushPlus
	TypeWePlus     ChannelType = "WePlus"     // 微加机器人
	TypeQMSG       ChannelType = "QMSG"       // Qmsg酱
	TypeWeCom      ChannelType = "WeCom"      // 企业微信
	TypeTelegram   ChannelType = "Telegram"   // Telegram
	TypeAIBOTK     ChannelType = "AIBOTK"     // 智能微秘书
	TypeEmail      ChannelType = "Email"      // 邮箱
	TypePushMe     ChannelType = "PushMe"     // PushMe
	TypeChronocat  ChannelType = "Chronocat"  // Chronocat
	TypeWebhook    ChannelType = "Webhook"    // Webhook
	TypeNtfy       ChannelType = "Ntfy"       // Ntfy
	TypeWxPusher   ChannelType = "WxPusher"   // WxPusher
)