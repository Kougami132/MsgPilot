package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/config"
	"github.com/kougami132/MsgPilot/internal/service"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
	env         *config.Env
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService, env *config.Env) *AuthController {
	return &AuthController{
		authService: authService,
		env:         env,
	}
}

// RegisterRoutes 注册路由
func (c *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/register", c.Register)
		auth.POST("/refresh", c.RefreshToken)
		auth.POST("/changePassword", c.ChangePassword)
	}
}

// Login godoc
// @Summary 登录
// @Description 登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {object} object "ok"
// @Failure 400 {object} object "无效的输入或验证错误"
// @Failure 401 {object} object "未授权"
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResponse, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse)
}

// Register godoc
// @Summary 注册
// @Description 注册
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {object} object "ok"
// @Failure 400 {object} object "无效的输入或验证错误"
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenResponse, err := c.authService.Register(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, tokenResponse)
}

// RefreshToken godoc
// @Summary 刷新令牌
// @Description 刷新令牌
// @Tags Auth
// @Accept json
// @Produce json
// @Param token body string true "Token"
// @Success 200 {object} object "ok"
// @Failure 400 {object} object "无效的输入或验证错误"
// @Failure 401 {object} object "未授权"
// @Router /api/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	type RefreshRequest struct {
		Token string `json:"token" binding:"required"`
	}

	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, expiry, err := c.authService.RefreshToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expiry":       expiry,
	})
}

// ChangePassword godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param old_password body string true "Old Password"
// @Param new_password body string true "New Password"
// @Success 200 {object} object "ok"
// @Failure 400 {object} object "无效的输入或验证错误"
// @Router /api/auth/changePassword [post]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.authService.ChangePassword(req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
