package test

import (
	"fmt"

	"github.com/kougami132/MsgPilot/channels"
	"github.com/kougami132/MsgPilot/models"
	"gorm.io/datatypes"
)

func Test() {
	testChannel := models.Channel {
		ID: "1",
		Name: "test",
		Type: "onebot",
		Direction: "out",
		Config: datatypes.JSON{},
	}
	onebot, err := channels.GetChannelType(testChannel)
	if err != nil {
		fmt.Println(err)
	}
	onebot.Send("test")
}
