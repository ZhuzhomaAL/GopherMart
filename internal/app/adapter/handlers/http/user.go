package http

import (
	"encoding/json"
	"errors"
	domenuser "github.com/ZhuzhomaAL/GopherMart/internal/app/core/domain/user"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/core/ports/service"
	"github.com/ZhuzhomaAL/GopherMart/internal/app/infra/auth"
	"io"
	"net/http"
)

type UserHandler struct {
	us service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{us: us}
}

type userCreds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	creds := userCreds{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		if err == io.EOF {
			http.Error(w, "request is empty, expected not empty", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	user, err := u.us.Register(r.Context(), creds.Login, creds.Password)
	if err != nil {
		if errors.Is(err, domenuser.LoginAlreadyExists{}) {
			http.Error(w, "internal server error occurred", http.StatusConflict)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
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
	return
}

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	creds := userCreds{}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		if err == io.EOF {
			http.Error(w, "incorrect response", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	user, err := u.us.Login(r.Context(), creds.Login, creds.Password)
	if err != nil {
		if errors.Is(err, domenuser.IncorrectLoginOrPassword{}) {
			http.Error(w, "internal server error occurred", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "internal server error occurred", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	http.SetCookie(
		w, &http.Cookie{
			Name:  "token",
			Value: tokenString,
			Path:  "/",
		},
	)
	return
}