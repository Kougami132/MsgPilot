package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/internal/service"
	"github.com/kougami132/MsgPilot/internal/channels"
)

type AdapterController struct {
	handlerService service.HandlerService
}

func NewAdapterController(handlerService service.HandlerService) *AdapterController {
	return &AdapterController{handlerService: handlerService}
}

func (c *AdapterController) RegisterRoutes(router *gin.RouterGroup) {
	adapterRoutes := router.Group("/adapter")
	{
		adapterRoutes.GET("/list", c.GetChannels)
	}
	onebotRoutes := router.Group("/onebot")
	{
		onebotRoutes.GET("/:ticket/send_msg", c.OneBotSendMsg)
		onebotRoutes.POST("/:ticket/send_msg", c.OneBotSendMsg)
	}
	barkRoutes := router.Group("/bark")
	{
		barkRoutes.GET("/:ticket/*action", c.BarkSendMsg)
		barkRoutes.POST("/:ticket", c.BarkSendMsg)
		barkRoutes.POST("/push", c.BarkSendMsg)
	}
	gotifyRoutes := router.Group("/gotify")
	{
		gotifyRoutes.POST("/message", c.GotifySendMsg)
	}
	pushdeerRoutes := router.Group("/pushdeer")
	{
		pushdeerRoutes.POST("/message/push", c.PushDeerSendMsg)
	}
	ntfyRoutes := router.Group("/ntfy")
	{
		ntfyRoutes.POST("/:ticket", c.NtfySendMsg)
	}
	webhookRoutes := router.Group("/webhook")
	{
		webhookRoutes.GET("/:ticket", c.WebhookSendMsg)
		webhookRoutes.POST("/:ticket", c.WebhookSendMsg)
	}
}

// GetChannels godoc
// @Summary 获取所有适配器列表
// @Description 获取所有适配器列表
// @Tags Adapter
// @Accept json
// @Produce json
// @Success 200 {object} object "ok"
// @Router /api/adapter/list [get]
func (c *AdapterController) GetChannels(ctx *gin.Context) {
	adapters := channels.GetChannelAdapters()
	handlers := channels.GetChannelHandlers()
	ctx.JSON(http.StatusOK, gin.H{
		"adapters": adapters,
		"handlers": handlers,
	})
}

// OneBotSendMsg godoc
// @Summary 接收OneBot消息
// @Description 接收OneBot消息
// @Tags Adapter
// @Accept json
// @Produce json
// @Param ticket path string true "Ticket"
// @Param message formData string false "Message"
// @Success 200 {object} object "ok"
// @Router /api/onebot/{ticket}/send_msg [post]
func (c *AdapterController) OneBotSendMsg(ctx *gin.Context) {
	ticket := ctx.Param("ticket")

	// 尝试从不同来源获取消息内容
	var msg string
	
	if ctx.Request.Method == "GET" {
		// 从Query参数获取
		msg = ctx.Query("message")
	} else if ctx.Request.Method == "POST" {
		// 从POST form获取
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
	}

	// 调用service处理消息
	message, err := c.handlerService.OneBotPush(ticket, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

// BarkSendMsg godoc
// @Summary 接收Bark消息
// @Description 接收Bark消息
// @Tags Adapter
// @Accept json
// @Produce json
// @Param ticket path string true "Ticket"
// @Param title formData string false "Title"
// @Param subtitle formData string false "Subtitle"
// @Param body formData string false "Body"
// @Success 200 {object} object "ok"
// @Router /api/bark/{ticket}/push [post]
func (c *AdapterController) BarkSendMsg(ctx *gin.Context) {
	ticket := ctx.Param("ticket")

	var title, subtitle, body string

	if ctx.Request.Method == "GET" {
		action := ctx.Param("action") //  例如 "/body" 或 "/title/body" 或 "/title/subtitle/body"
		trimmedAction := strings.TrimPrefix(action, "/")
		parts := strings.Split(trimmedAction, "/")

		switch len(parts) {
		case 1: // /:body
			if parts[0] != "" { // 确保URL不是 /:key/ 这样的形式
				body = parts[0]
			}
		case 2: // /:title/:body
			title = parts[0]
			body = parts[1]
		case 3: // /:title/:subtitle/:body
			title = parts[0]
			subtitle = parts[1]
			body = parts[2]
		default:
			break
		}
	} else if ctx.Request.Method == "POST" {
		// 从POST form获取
		if title == "" && body == "" {
			title = ctx.PostForm("title")
			subtitle = ctx.PostForm("subtitle")
			body = ctx.PostForm("body")
		}

		// 从JSON body获取
		if title == "" && body == "" {
			var jsonData struct {
				Title string `json:"title"`
				Subtitle string `json:"subtitle"`
				Body  string `json:"body"`
				DeviceKey string `json:"device_key"`
			}
			if err := ctx.ShouldBindJSON(&jsonData); err == nil {
				if title == "" {
					title = jsonData.Title
				}
				if subtitle == "" {
					subtitle = jsonData.Subtitle
				}
				if body == "" {
					body = jsonData.Body
				}
				if ticket == "" {
					ticket = jsonData.DeviceKey
				}
			}
		}
	}

	// 副标题合并到标题中
	if subtitle != "" {
		title = title + " - " + subtitle
	}

	message, err := c.handlerService.CommonPush(ticket, title, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

// GotifySendMsg godoc
// @Summary 接收Gotify消息
// @Description 接收Gotify消息
// @Tags Adapter
// @Accept json
// @Produce json
// @Param token query string true "Token"
// @Param title formData string false "Title"
// @Param message formData string false "Message"
// @Success 200 {object} object "ok"
// @Router /api/gotify/message [post]
func (c *AdapterController) GotifySendMsg(ctx *gin.Context) {
	ticket := ctx.Query("token")

	title := ctx.PostForm("title")
	msg := ctx.PostForm("message")

	message, err := c.handlerService.CommonPush(ticket, title, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

// PushDeerSendMsg godoc
// @Summary 接收PushDeer消息
// @Description 接收PushDeer消息
// @Tags Adapter
// @Accept json
// @Produce json
// @Param token formData string true "Token"
// @Param text formData string false "Text"
// @Param desp formData string false "Desp"
// @Success 200 {object} object "ok"
// @Router /api/pushdeer/message/push [post]
func (c *AdapterController) PushDeerSendMsg(ctx *gin.Context) {
	ticket := ctx.PostForm("token")
	title := ctx.PostForm("text")
	msg := ctx.PostForm("desp")

	message, err := c.handlerService.CommonPush(ticket, title, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

// NtfySendMsg godoc
// @Summary 接收Ntfy消息
// @Description 接收Ntfy消息
// @Tags Adapter
// @Accept json
// @Produce json
// @Param ticket path string true "Ticket"
// @Param title formData string false "Title"
// @Param message formData string false "Message"
// @Success 200 {object} object "ok"
// @Router /api/ntfy/{ticket} [post]
func (c *AdapterController) NtfySendMsg(ctx *gin.Context) {
	ticket := ctx.Param("ticket")
	title := ctx.GetHeader("Title")
	if title == "" {
		title = ctx.PostForm("title")
	}
	msg := ctx.PostForm("message")

	if ticket == "" && title == "" && msg == "" {
		var jsonData struct {
			Topic   string `json:"topic"`
			Title   string `json:"title"`
			Message string `json:"message"`
		}
		if err := ctx.ShouldBindJSON(&jsonData); err == nil {
			ticket = jsonData.Topic
			title = jsonData.Title
			msg = jsonData.Message
		}
	}

	message, err := c.handlerService.CommonPush(ticket, title, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

// WebhookSendMsg godoc
// @Summary 接收Webhook消息
// @Description 接收Webhook消息
// @Tags Adapter
// @Accept json
// @Produce json
// @Param ticket path string true "Ticket"
// @Param title formData string false "Title"
// @Param message formData string false "Message"
// @Success 200 {object} object "ok"
// @Router /api/webhook/{ticket} [post]
func (c *AdapterController) WebhookSendMsg(ctx *gin.Context) {
	ticket := ctx.Param("ticket")
	title := ctx.Query("title")
	if title == "" {
		title = ctx.PostForm("title")
	}
	msg := ctx.Query("message")
	if msg == "" {
		msg = ctx.PostForm("message")
	}

	message, err := c.handlerService.CommonPush(ticket, title, msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}
