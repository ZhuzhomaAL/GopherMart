package postgres

import (
	"context"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/storage/postgres"
)

type UserRepository struct {
	client *postgres.Client
}

func NewUserRepository(client *postgres.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (ur UserRepository) CreateUser(ctx context.Context, user user.User) (user.User, error) {
	_, err := ur.client.NewInsert().Model(&user).Exec(ctx)
	return user, err
}

func (ur UserRepository) GetByLogin(ctx context.Context, login string) (user.User, error) {
	u := new(user.User)
	err := ur.client.NewSelect().Model(u).Where("login = ?", login).Scan(ctx)
	return *u, err
}
