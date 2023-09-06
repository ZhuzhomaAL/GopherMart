package service

import (
	"context"
	"errors"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/adapters/repository/mocks"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func generateLogin() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TestUserService_Login(t *testing.T) {
	ctx := context.Background()
	ID, _ := uuid.NewV7()
	login := generateLogin()
	password := strconv.Itoa(rand.New(rand.NewSource(10)).Int())
	passwordHash, _ := auth.MakePasswordHash(password)
	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		mockRes user.User
		wantErr bool
		mockErr error
		expErr  error
	}{
		{
			name: "Test_1.Проверка успешной авторизации",
			args: args{
				ctx:      ctx,
				login:    login,
				password: password,
			},
			mockRes: user.User{
				ID:       ID,
				Login:    login,
				Password: passwordHash,
			},
			wantErr: false,
			mockErr: nil,
		},
		{
			name: "Test_2.Метод возвращает ошибку.Некорректный логин или пароль",
			args: args{
				ctx:      ctx,
				login:    login,
				password: password,
			},
			mockRes: user.User{
				ID:       ID,
				Login:    login,
				Password: passwordHash,
			},
			wantErr: true,
			mockErr: repository.NoResultError{},
			expErr: &user.IncorrectLoginOrPassword{
				Login:    login,
				Password: password,
			},
		},
		{
			name: "Test_3.Метод возвращает ошибку.Сравнение пароля и хеша",
			args: args{
				ctx:      ctx,
				login:    login,
				password: password,
			},
			mockRes: user.User{
				ID:       ID,
				Login:    login,
				Password: password,
			},
			wantErr: true,
			mockErr: nil,
			expErr: &user.IncorrectLoginOrPassword{
				Login:    login,
				Password: password,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.UserRepository{}
				us := NewUserService(&rep)
				rep.On("GetByLogin", tt.args.ctx, tt.args.login).Return(tt.mockRes, tt.mockErr)
				loginData, err := us.Login(tt.args.ctx, tt.args.login, tt.args.password)
				if (err != nil) != tt.wantErr {
					t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					require.Equal(t, tt.mockRes, loginData)
				}
				if tt.wantErr {
					require.Equal(t, tt.expErr, err)
				}
			},
		)
	}
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	ID, _ := uuid.NewV7()
	login := generateLogin()
	password := strconv.Itoa(rand.New(rand.NewSource(10)).Int())
	passwordHash, _ := auth.MakePasswordHash(password)
	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name          string
		args          args
		mockRes       user.User
		wantErr       bool
		mockErr       error
		mockCreateErr error
	}{
		{
			name: "Test_1.Проверка успешной регистрации",
			args: args{
				ctx:      ctx,
				login:    login,
				password: password,
			},
			mockRes: user.User{
				ID:       ID,
				Login:    login,
				Password: passwordHash,
			},
			wantErr: false,
			mockErr: errors.New("no data"),
		},
		{
			name: "Test_2.Метод возвращает ошибку.Логин уже существует",
			args: args{
				ctx:      ctx,
				login:    login,
				password: password,
			},
			mockRes: user.User{
				ID:       ID,
				Login:    login,
				Password: passwordHash,
			},
			wantErr: true,
			mockErr: nil,
		},
		{
			name: "Test_3.Метод возвращает ошибку.Не удалось создать пользователя",
			args: args{
				ctx:      ctx,
				login:    login,
				password: password,
			},
			mockRes: user.User{
				ID:       ID,
				Login:    login,
				Password: passwordHash,
			},
			wantErr:       true,
			mockErr:       errors.New("no data"),
			mockCreateErr: errors.New("can't create"),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				rep := mocks.UserRepository{}
				us := NewUserService(&rep)
				rep.On("GetByLogin", tt.args.ctx, tt.args.login).Return(tt.mockRes, tt.mockErr)
				rep.On("CreateUser", tt.args.ctx, mock.AnythingOfType("user.User")).Return(tt.mockRes, tt.mockCreateErr)
				userData, err := us.Register(tt.args.ctx, tt.args.login, tt.args.password)
				if (err != nil) != tt.wantErr {
					t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr {
					require.Equal(t, tt.mockRes, userData)
				}
			},
		)
	}
}
