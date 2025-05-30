package usecase

import (
	"errors"
	"time"

	"github.com/kougami132/MsgPilot/config"
	"github.com/kougami132/MsgPilot/internal/utils"
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
	"golang.org/x/crypto/bcrypt"
)

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

// AuthUseCase 认证用例接口
type AuthUseCase interface {
	Login(username, password string) (*TokenResponse, error)
	Register(username, password string) (*TokenResponse, error)
	RefreshToken(refreshToken string) (string, int64, error)
	ChangePassword(username, oldPassword, newPassword string) error
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

	expiryTimestamp := time.Now().Unix() + int64(u.env.AccessTokenExpiryHour*3600)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expiryTimestamp,
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

	expiryTimestamp := time.Now().Unix() + int64(u.env.AccessTokenExpiryHour*3600)

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       expiryTimestamp,
	}, nil
}

// RefreshToken 刷新令牌
func (u *AuthUseCaseImpl) RefreshToken(refreshToken string) (string, int64, error) {
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

	expiryTimestamp := time.Now().Unix() + int64(u.env.AccessTokenExpiryHour*3600)
	return accessToken, expiryTimestamp, nil
}

// ChangePassword 修改密码
func (u *AuthUseCaseImpl) ChangePassword(username, oldPassword, newPassword string) error {
	// 查找用户
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		// 为了安全，不明确指出是用户不存在还是密码错误
		return errors.New("用户名或旧密码错误")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("用户名或旧密码错误")
	}

	// 检查新旧密码是否相同
	if oldPassword == newPassword {
		return errors.New("新密码不能与旧密码相同")
	}

	// 加密新密码
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("新密码加密失败")
	}

	// 更新密码
	user.Password = string(hashedNewPassword)
	if err := u.userRepo.Update(user); err != nil {
		return errors.New("更新密码失败")
	}

	return nil
}
