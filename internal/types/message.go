package types

// MessageStatus 消息状态枚举定义
type MessageStatus int
const (
    StatusPending MessageStatus = iota // 待发送
    StatusSending                      // 发送中
    StatusSuccess                      // 已发送
    StatusFailed                       // 失败
)