package repository

import (
	"errors"
	"fmt"

	"github.com/fajritsaniy/golang-SHM/model"
	"github.com/fajritsaniy/golang-SHM/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository[model.UserCredential]
	GetByUsername(username string) (*model.UserCredential, error)
	GetByUsernamePassword(username string, password string) (*model.UserCredential, error)
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) Search(by map[string]interface{}) ([]model.UserCredential, error) {
	var users []model.UserCredential
	result := u.db.Where(by).Find(&users).Error
	if result != nil {
		return nil, result
	}
	return users, nil
}

func (u *userRepository) List() ([]model.UserCredential, error) {
	var users []model.UserCredential
	result := u.db.Find(&users).Error
	if result != nil {
		return nil, result
	}
	return users, nil
}

func (u *userRepository) Get(id string) (*model.UserCredential, error) {
	var user model.UserCredential
	result := u.db.First(&user, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &user, nil
}

func (u *userRepository) Save(payload *model.UserCredential) error {
	return u.db.Save(payload).Error
}

func (u *userRepository) Delete(id string) error {
	return u.db.Delete(&model.UserCredential{}, "id=?", id).Error
}

func (u *userRepository) GetByUsernameActive(username string) (*model.UserCredential, error) {
	var userCredential model.UserCredential
	result := u.db.Where("user_name = ?", username).Where("is_active = ?", true).First(&userCredential)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username '%s' not found", username)
		}
		return nil, fmt.Errorf("failed to get user with username '%s': %v", username, err)
	}
	return &userCredential, nil
}

func (u *userRepository) GetByUsername(username string) (*model.UserCredential, error) {
	var userCredential model.UserCredential
	result := u.db.Where("user_name = ?", username).First(&userCredential)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username '%s' not found", username)
		}
		return nil, fmt.Errorf("failed to get user with username '%s': %v", username, err)
	}
	return &userCredential, nil
}

func (u *userRepository) GetByUsernamePassword(username string, password string) (*model.UserCredential, error) {
	user, err := u.GetByUsernameActive(username)
	if err != nil {
		return nil, err
	}

	pwdCheck := utils.CheckPasswordHash(password, user.Password)
	if !pwdCheck {
		return nil, fmt.Errorf("Password don't match")
	}
	return user, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
