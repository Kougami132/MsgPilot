package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/usecase"
)

type AdapterController struct {
	adapterUsecase usecase.AdapterUsecase
}

func NewAdapterController(adapterUsecase usecase.AdapterUsecase) *AdapterController {
	return &AdapterController{adapterUsecase: adapterUsecase}
}

func (c *AdapterController) RegisterRoutes(router *gin.RouterGroup) {
	adapterRoutes := router.Group("/adapter")
	{
		adapterRoutes.Any("/:ticket/send_msg", c.OneBotSendMsg)
	}
}

func (c *AdapterController) OneBotSendMsg(ctx *gin.Context) {
	ticket := ctx.Param("ticket")

	// 尝试从不同来源获取消息内容
	var msg string
	
	// 从Query参数获取
	msg = ctx.Query("message")
	
	// 如果Query为空,尝试从POST form获取
	if msg == "" {
		msg = ctx.PostForm("message")
	}
	
	// 如果form为空,尝试从JSON body获取
	if msg == "" {
		var jsonData struct {
			Message string `json:"message"`
		}
		if err := ctx.ShouldBindJSON(&jsonData); err == nil {
			msg = jsonData.Message
		}
	}

	// 调用usecase处理消息
	message, err := c.adapterUsecase.OneBotSendMessage(ticket, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}