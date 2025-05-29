package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/usecase"
)

type ConfigController struct {
	configUsecase usecase.ConfigUsecase
}

func NewConfigController(configUsecase usecase.ConfigUsecase) *ConfigController {
	return &ConfigController{configUsecase: configUsecase}
}

func (c *ConfigController) RegisterRoutes(router *gin.RouterGroup) {
	configRoutes := router.Group("/config")
	{
		configRoutes.POST("/create", c.CreateConfig)
		configRoutes.GET("/list", c.GetAllConfigs)
		configRoutes.GET("/get/:key", c.GetConfigByKey)
		configRoutes.PUT("/update/:key", c.UpdateConfig)
		configRoutes.DELETE("/delete/:key", c.DeleteConfig)
	}
}

func (c *ConfigController) CreateConfig(ctx *gin.Context) {
	var config models.Config
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.configUsecase.CreateConfig(&config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, config)
}

func (c *ConfigController) GetAllConfigs(ctx *gin.Context) {
	configs, err := c.configUsecase.GetAllConfigs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, configs)
}

func (c *ConfigController) GetConfigByKey(ctx *gin.Context) {
	key := ctx.Param("key")
	config, err := c.configUsecase.GetConfigByKey(key)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

func (c *ConfigController) UpdateConfig(ctx *gin.Context) {
	key := ctx.Param("key")
	var config models.Config
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.Key = key // Ensure Key is set from path parameter
	if err := c.configUsecase.UpdateConfig(&config); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

func (c *ConfigController) DeleteConfig(ctx *gin.Context) {
	key := ctx.Param("key")
	if err := c.configUsecase.DeleteConfig(key); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Config deleted successfully"})
}
