package service

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/errors"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"log/slog"
)

type UserService struct {
	repo repository.Repository
	log  *slog.Logger
}

func NewUserService(log *slog.Logger, repo repository.Repository) *UserService {
	return &UserService{log: log, repo: repo}
}

func (s *UserService) Register(ctx context.Context, login string, password string) (token string, user *models.User, err error) {
	const op = "UserService.Register"
	log := s.log.With(
		slog.String("op", op),
	)

	user = &models.User{}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Error(err.Error())
		return token, user, err
	}

	user.Username = login
	user.Password = hashedPassword

	user, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		if !cstErrors.IsCustomError(err) {
			log.Error(err.Error())
		}
		return token, nil, err
	}

	token, err = utils.GenerateJWT(user.ID)
	if err != nil {
		log.Error(err.Error())
	}
	return token, user, nil
}

func (s *UserService) Login(ctx context.Context, login string, password string) (token string, user *models.User, err error) {
	const op = "UserService.Login"
	log := s.log.With(
		slog.String("op", op),
	)

	user = &models.User{}

	user, err = s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return token, nil, cstErrors.InvalidCredentials
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return token, nil, cstErrors.InvalidCredentials
	}

	go s.repo.FixLastActivity(ctx, user.ID)

	token, err = utils.GenerateJWT(user.ID)
	if err != nil {
		log.Error(err.Error())
	}
	return token, user, err
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (user *models.User, err error) {
	const op = "UserService.GetUserByID"
	log := s.log.With(
		slog.String("op", op),
	)

	user, err = s.repo.GetUserByID(ctx, id)
	if err != nil && !cstErrors.IsCustomError(err) {
		log.Error(err.Error())
	}
	return user, err
}
