package channels

type ChannelHandler interface {
    Send(content string) error
}