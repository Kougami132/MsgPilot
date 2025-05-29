package usecase

import (
	"errors"

	"github.com/kougami132/MsgPilot/config"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
	"github.com/kougami132/MsgPilot/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthUseCase 认证用例接口
type AuthUseCase interface {
	Login(username, password string) (*TokenResponse, error)
	Register(username, password string) (*TokenResponse, error)
	RefreshToken(refreshToken string) (string, int, error)
}

// AuthUseCaseImpl 认证用例实现
type AuthUseCaseImpl struct {
	userRepo repository.UserRepository
	env      *config.Env
}

// NewAuthUseCase 创建认证用例
func NewAuthUseCase(userRepo repository.UserRepository, env *config.Env) AuthUseCase {
	return &AuthUseCaseImpl{
		userRepo: userRepo,
		env:      env,
	}
}

// Login 登录
func (u *AuthUseCaseImpl) Login(username, password string) (*TokenResponse, error) {
	// 查找用户
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成令牌
	accessToken, err := utils.GenerateToken(user.Username, u.env.AccessTokenSecret, u.env.AccessTokenExpiryHour)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	refreshToken, err := utils.GenerateToken(user.Username, u.env.RefreshTokenSecret, u.env.RefreshTokenExpiryHour)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    u.env.AccessTokenExpiryHour * 3600,
	}, nil
}

// Register 注册
func (u *AuthUseCaseImpl) Register(username, password string) (*TokenResponse, error) {
	// 检查用户名是否已存在
	_, err := u.userRepo.GetByUsername(username)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.Create(&user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	// 生成令牌
	accessToken, err := utils.GenerateToken(user.Username, u.env.AccessTokenSecret, u.env.AccessTokenExpiryHour)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	refreshToken, err := utils.GenerateToken(user.Username, u.env.RefreshTokenSecret, u.env.RefreshTokenExpiryHour)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    u.env.AccessTokenExpiryHour * 3600,
	}, nil
}

// RefreshToken 刷新令牌
func (u *AuthUseCaseImpl) RefreshToken(refreshToken string) (string, int, error) {
	// 验证刷新令牌
	claims, err := utils.ValidateToken(refreshToken, u.env.RefreshTokenSecret)
	if err != nil {
		return "", 0, errors.New("无效的刷新令牌")
	}

	// 生成新的访问令牌
	accessToken, err := utils.GenerateToken(claims.Username, u.env.AccessTokenSecret, u.env.AccessTokenExpiryHour)
	if err != nil {
		return "", 0, errors.New("生成令牌失败")
	}

	return accessToken, u.env.AccessTokenExpiryHour * 3600, nil
}
