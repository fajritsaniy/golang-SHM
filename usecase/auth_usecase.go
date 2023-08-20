package usecase

import (
	"errors"
	"fmt"

	"github.com/fajritsaniy/golang-SHM/model"
	"github.com/fajritsaniy/golang-SHM/utils"
	"github.com/fajritsaniy/golang-SHM/utils/security"
	"gorm.io/gorm"

	"github.com/fajritsaniy/golang-SHM/repository"
)

type AuthenticationUseCase interface {
	Login(username string, password string) (string, error)
	Register(payload *model.UserCredential) error
	UserActivation(payload *model.UserCredential) (bool, error)
}

type authenticationUseCase struct {
	repo         repository.UserRepository
	tokenService security.AccessToken
}

func (a *authenticationUseCase) Login(username string, password string) (string, error) {
	user, err := a.repo.GetByUsernamePassword(username, password)
	var token string
	if err != nil {
		return "", fmt.Errorf("user with username: %s not found", username)
	}
	if user != nil {
		token, err = a.tokenService.CreateAccessToken(user)
		fmt.Println("token:", token)
		if err != nil {
			return "", err
		}
	}
	return token, nil
}

func (a *authenticationUseCase) Register(payload *model.UserCredential) error {
	if payload.ID != "" {
		_, err := a.repo.Get(payload.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user with ID '%s' not found", payload.ID)
			}
			return fmt.Errorf("failed to check user with ID '%s': %v", payload.ID, err)
		}
	}

	if payload.Password != "" {
		password, err := utils.HashPassword(payload.Password)
		if err != nil {
			return err
		}
		payload.Password = password
	}

	return a.repo.Save(payload)
}

func (a *authenticationUseCase) UserActivation(payload *model.UserCredential) (bool, error) {
	user, err := a.repo.GetByUsername(payload.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("username '%s' not found", payload.UserName)
		}
		return false, fmt.Errorf("failed to check user with username '%s': %v", payload.UserName, err)
	}

	var status bool
	if user.IsActive {
		user.IsActive = false
	} else {
		user.IsActive = true
	}
	status = user.IsActive

	return status, a.repo.Save(user)
}

func NewAuthenticationUseCase(repo repository.UserRepository, tokenService security.AccessToken) AuthenticationUseCase {
	return &authenticationUseCase{repo: repo, tokenService: tokenService}
}
