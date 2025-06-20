package main

import (
	"fmt"

	"github.com/kougami132/MsgPilot/internal/channels"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/internal/types"
	"gorm.io/datatypes"
)

func Send() {
	testChannel := models.Channel{
		ID:     1,
		Name:   "test",
		Type:   types.TypeOneBot,
		Config: datatypes.JSON(`{"base_url":"https://bot.kougami.me/send_msg","message_type":"private","user_id":"1329623049"}`),
	}
	onebot, err := channels.GetChannelHandler(testChannel)
	if err != nil {
		fmt.Println(err)
	}
	onebot.Send(&models.Message{
		Title:   "test",
		Content: "test",
	})
}

func main() {
	Send()
}
