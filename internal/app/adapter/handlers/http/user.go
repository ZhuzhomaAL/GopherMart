package http

import (
	"encoding/json"
	domenuser "github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/logger"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	us  service.UserService
	log logger.MyLogger
}

func NewUserHandler(us service.UserService, log logger.MyLogger) *UserHandler {
	return &UserHandler{us: us, log: log}
}

type userCreds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	creds := userCreds{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		u.log.L.Error("failed to decode request", zap.Error(err))
		http.Error(w, "Can not parse request", http.StatusBadRequest)
		return
	}

	user, err := u.us.Register(r.Context(), creds.Login, creds.Password)
	if err != nil {
		u.log.L.Error("failed to register user", zap.Error(err))
		if _, ok := err.(*domenuser.LoginAlreadyExists); ok {
			http.Error(w, "internal server error occurred", http.StatusConflict)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		u.log.L.Error("failed to generate auth token", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
	}
	http.SetCookie(
		w, cookie,
	)
	w.WriteHeader(http.StatusOK)
}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	creds := userCreds{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		u.log.L.Error("failed to decode request", zap.Error(err))
		http.Error(w, "Can not parse request", http.StatusBadRequest)
		return
	}

	user, err := u.us.Login(r.Context(), creds.Login, creds.Password)
	if err != nil {
		u.log.L.Error("failed to login user", zap.Error(err))
		if _, ok := err.(*domenuser.IncorrectLoginOrPassword); ok {
			http.Error(w, "internal server error occurred", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		u.log.L.Error("failed to generate auth token", zap.Error(err))
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
	}
	http.SetCookie(
		w, cookie,
	)
	w.WriteHeader(http.StatusOK)
}
