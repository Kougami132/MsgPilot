package controller

import (
	"net/http"
	"strconv"
	
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

// CreateMessage godoc
// @Summary 创建消息
// @Description 创建一个新的消息
// @Tags Message
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param message body models.Message true "Message"
// @Success 201 {object} models.Message
// @Failure 400 {object} map[string]string "无效的输入或验证错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/message/create [post]
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

// GetAllMessages godoc
// @Summary 获取所有消息
// @Description 获取所有消息的列表
// @Tags Message
// @Produce json
// @Param Authorization header string true "Authorization"
// @Success 200 {array} models.Message
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/message/list [get]
func (c *MessageController) GetAllMessages(ctx *gin.Context) {
	messages, err := c.messageUsecase.GetAllMessages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, messages)
}

// GetMessageByID godoc
// @Summary 根据ID获取消息
// @Description 根据提供的ID获取单个消息的详细信息
// @Tags Message
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path string true "Message ID"
// @Success 200 {object} models.Message
// @Failure 404 {object} map[string]string "消息未找到"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/message/get/{id} [get]
func (c *MessageController) GetMessageByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	message, err := c.messageUsecase.GetMessageByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

// UpdateMessage godoc
// @Summary 更新消息
// @Description 根据提供的ID更新现有的消息
// @Tags Message
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path string true "Message ID"
// @Param message body models.Message true "Message"
// @Success 200 {object} models.Message
// @Failure 400 {object} map[string]string "无效的输入或验证错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/message/update/{id} [put]
func (c *MessageController) UpdateMessage(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var message models.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message.ID = id
	if err := c.messageUsecase.UpdateMessage(&message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, message)
}

// DeleteMessage godoc
// @Summary 删除消息
// @Description 根据提供的ID删除消息
// @Tags Message
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path string true "Message ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string "消息未找到"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/message/delete/{id} [delete]
func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.messageUsecase.DeleteMessage(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
