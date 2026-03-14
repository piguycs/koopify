package auth

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type KoopifyClaims struct {
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
}

var (
	ErrMissingToken   = errors.New("missing auth token")
	ErrInvalidToken   = errors.New("invalid auth token")
	ErrInvalidSubject = errors.New("invalid token subject")
)

func ClaimsFromContext(ctx *echo.Context) (*KoopifyClaims, error) {
	user := ctx.Get("user")
	if user == nil {
		return nil, ErrMissingToken
	}

	token, ok := user.(*jwt.Token)
	if !ok || token == nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*KoopifyClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func UserIDFromToken(ctx *echo.Context) (int64, error) {
	claims, err := ClaimsFromContext(ctx)
	if err != nil {
		return 0, err
	}

	if claims.Subject == "" {
		return 0, ErrInvalidSubject
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, ErrInvalidSubject
	}

	return userID, nil
}

func IsAdminFromToken(ctx *echo.Context) bool {
	claims, err := ClaimsFromContext(ctx)
	if err != nil {
		return false
	}

	return claims.Admin
}
