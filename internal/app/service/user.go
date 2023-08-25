package service

import (
	"context"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository"
	"github.com/gofrs/uuid"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us UserService) Register(ctx context.Context, login, password string) (user.User, error) {
	if _, err := us.repo.GetByLogin(ctx, login); err == nil {
		return user.User{}, &user.LoginAlreadyExists{Login: login}
	}
	passwordHash, err := user.MakePasswordHash(password)
	if err != nil {
		return user.User{}, err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return user.User{}, err
	}

	return us.repo.CreateUser(
		ctx, user.User{
			ID:       id,
			Login:    login,
			Password: passwordHash,
		},
	)
}

func (us UserService) Login(ctx context.Context, login, password string) (user.User, error) {
	u, err := us.repo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.NoResultError{}) {
			return user.User{}, &user.IncorrectLoginOrPassword{
				Login:    login,
				Password: password,
			}
		}
		return user.User{}, err
	}

	if !u.CheckPassword(password) {
		return user.User{}, &user.IncorrectLoginOrPassword{
			Login:    login,
			Password: password,
		}
	}

	return u, nil
}
