package service

import (
	"github.com/kougami132/MsgPilot/internal/repository"
	"github.com/kougami132/MsgPilot/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	CheckPassword(username, password string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (u *userService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return u.userRepo.Create(user)
}

func (u *userService) GetAllUsers() ([]models.User, error) {
	return u.userRepo.GetAll()
}

func (u *userService) GetUserByID(id int) (*models.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userService) GetUserByUsername(username string) (*models.User, error) {
	return u.userRepo.GetByUsername(username)
}

func (u *userService) UpdateUser(user *models.User) error {
	// 如果密码字段不为空，则哈希密码
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return u.userRepo.Update(user)
}

func (u *userService) DeleteUser(id int) error {
	return u.userRepo.Delete(id)
}

func (u *userService) CheckPassword(username, password string) (*models.User, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // 密码不匹配
	}
	return user, nil
}
