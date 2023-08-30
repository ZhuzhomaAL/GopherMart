package auth

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type contextUserIDKey int

const ContextUserID contextUserIDKey = iota

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

var jwtKey = []byte("gopher_mart")

func GenerateJWT(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{},
			UserID:           id,
		},
	)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var id uuid.UUID
			c, err := r.Cookie("token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			id, err = getUserIDFromToken(c.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserID, id)
			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func getUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtKey, nil
		},
	)
	if err != nil {
		return uuid.Nil, nil
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(ContextUserID).(uuid.UUID)
	return userID, ok
}
