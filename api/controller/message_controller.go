package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/usecase"
)

type MessageController struct {
	messageUsecase usecase.MessageUsecase
}

func NewMessageController(messageUsecase usecase.MessageUsecase) *MessageController {
	return &MessageController{messageUsecase: messageUsecase}
}

func (c *MessageController) RegisterRoutes(router *gin.RouterGroup) {
	messageRoutes := router.Group("/message")
	{
		messageRoutes.POST("/create", c.CreateMessage)
		messageRoutes.GET("/list", c.GetAllMessages)
		messageRoutes.GET("/get/:id", c.GetMessageByID)
		messageRoutes.PUT("/update/:id", c.UpdateMessage)
		messageRoutes.DELETE("/delete/:id", c.DeleteMessage)
	}
}

func (c *MessageController) CreateMessage(ctx *gin.Context) {
	var message models.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.messageUsecase.CreateMessage(&message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, message)
}

func (c *MessageController) GetAllMessages(ctx *gin.Context) {
	messages, err := c.messageUsecase.GetAllMessages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, messages)
}

func (c *MessageController) GetMessageByID(ctx *gin.Context) {
	id := ctx.Param("id")
	message, err := c.messageUsecase.GetMessageByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

func (c *MessageController) UpdateMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	var message models.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message.ID = id // Ensure ID is set from path parameter
	if err := c.messageUsecase.UpdateMessage(&message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.messageUsecase.DeleteMessage(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
