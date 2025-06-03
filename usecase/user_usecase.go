package usecase

import (
	"github.com/kougami132/MsgPilot/models"
	"github.com/kougami132/MsgPilot/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	CheckPassword(username, password string) (*models.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return u.userRepo.Create(user)
}

func (u *userUsecase) GetAllUsers() ([]models.User, error) {
	return u.userRepo.GetAll()
}

func (u *userUsecase) GetUserByID(id int) (*models.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *userUsecase) GetUserByUsername(username string) (*models.User, error) {
	return u.userRepo.GetByUsername(username)
}

func (u *userUsecase) UpdateUser(user *models.User) error {
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

func (u *userUsecase) DeleteUser(id int) error {
	return u.userRepo.Delete(id)
}

func (u *userUsecase) CheckPassword(username, password string) (*models.User, error) {
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
