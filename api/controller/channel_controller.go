package controller

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/usecase"
)

type ChannelController struct {
	channelUsecase usecase.ChannelUsecase
}

func NewChannelController(channelUsecase usecase.ChannelUsecase) *ChannelController {
	return &ChannelController{channelUsecase: channelUsecase}
}

func (c *ChannelController) RegisterRoutes(router *gin.RouterGroup) {
	channelRoutes := router.Group("/channel")
	{
		channelRoutes.POST("/create", c.CreateChannel)
		channelRoutes.GET("/list", c.GetAllChannels)
		channelRoutes.GET("/get/:id", c.GetChannelByID)
		channelRoutes.PUT("/update/:id", c.UpdateChannel)
		channelRoutes.DELETE("/delete/:id", c.DeleteChannel)
		channelRoutes.POST("/test", c.TestPush)
	}
}

// CreateChannel godoc
// @Summary 创建渠道
// @Description 创建一个新的渠道
// @Tags Channel
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param channel body models.Channel true "Channel"
// @Success 201 {object} models.Channel
// @Failure 400 {object} map[string]string "无效的输入或验证错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/channel/create [post]
func (c *ChannelController) CreateChannel(ctx *gin.Context) {
	var channel models.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.channelUsecase.CreateChannel(&channel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, channel)
}

// GetAllChannels godoc
// @Summary 获取所有渠道
// @Description 获取所有渠道的列表
// @Tags Channel
// @Produce json
// @Param Authorization header string true "Authorization"
// @Success 200 {array} models.Channel
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/channel/list [get]
func (c *ChannelController) GetAllChannels(ctx *gin.Context) {
	channels, err := c.channelUsecase.GetAllChannels()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, channels)
}

// GetChannelByID godoc
// @Summary 根据ID获取渠道
// @Description 根据提供的ID获取单个渠道的详细信息
// @Tags Channel
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path string true "Channel ID"
// @Success 200 {object} models.Channel
// @Failure 404 {object} map[string]string "渠道未找到"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/channel/get/{id} [get]
func (c *ChannelController) GetChannelByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	channel, err := c.channelUsecase.GetChannelByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}
	ctx.JSON(http.StatusOK, channel)
}

// UpdateChannel godoc
// @Summary 更新渠道
// @Description 根据提供的ID更新现有的渠道
// @Tags Channel
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path string true "Channel ID"
// @Param channel body models.Channel true "Channel"
// @Success 200 {object} models.Channel
// @Failure 400 {object} map[string]string "无效的输入或验证错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/channel/update/{id} [put]
func (c *ChannelController) UpdateChannel(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var channel models.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel.ID = id
	if err := c.channelUsecase.UpdateChannel(&channel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, channel)
}

// DeleteChannel godoc
// @Summary 删除渠道
// @Description 根据提供的ID删除现有的渠道
// @Tags Channel
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path string true "Channel ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string "渠道未找到"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/channel/delete/{id} [delete]
func (c *ChannelController) DeleteChannel(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.channelUsecase.DeleteChannel(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
}

// TestPush godoc
// @Summary 测试推送
// @Description 测试推送
// @Tags Channel
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param channel body models.Channel true "Channel"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "无效的输入或验证错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/channel/test [post]
func (c *ChannelController) TestPush(ctx *gin.Context) {
	var channel models.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.channelUsecase.TestPush(channel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "测试成功"})
}
