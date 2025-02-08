package service

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
)

type UserService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, login string, password string) (token string, user *models.User, err error) {
	user = &models.User{}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return token, user, err
	}

	user.Username = login
	user.Password = hashedPassword

	user, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return token, user, err
	}

	token, err = utils.GenerateJWT(user.ID)
	return token, user, nil
}

func (s *UserService) Login(ctx context.Context, login string, password string) (token string, user *models.User, err error) {
	user = &models.User{}

	user, err = s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return token, user, cstErrors.InvalidCredentials
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return token, user, cstErrors.InvalidCredentials
	}

	token, err = utils.GenerateJWT(user.ID)
	return token, user, err
}
