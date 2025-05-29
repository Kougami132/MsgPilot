package controller

import (
	"net/http"

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
	}
}

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

func (c *ChannelController) GetAllChannels(ctx *gin.Context) {
	channels, err := c.channelUsecase.GetAllChannels()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, channels)
}

func (c *ChannelController) GetChannelByID(ctx *gin.Context) {
	id := ctx.Param("id")
	channel, err := c.channelUsecase.GetChannelByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}
	ctx.JSON(http.StatusOK, channel)
}

func (c *ChannelController) UpdateChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	var channel models.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	channel.ID = id // Ensure ID is set from path parameter
	if err := c.channelUsecase.UpdateChannel(&channel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, channel)
}

func (c *ChannelController) DeleteChannel(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.channelUsecase.DeleteChannel(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
}
