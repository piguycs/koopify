package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"piguy.nl/koopify/internal"
)

type UserController struct {
	jwtSecret string
	service   UserService
}

func NewUserController(jwtSecret string, service UserService) UserController {
	return UserController{jwtSecret: jwtSecret, service: service}
}

func (uc *UserController) CreateUser(ctx *echo.Context) (err error) {
	user, err := internal.BindAndValidate[CreateUserRequest](ctx)

	if err != nil {
		return err
	}

	resp, err := uc.service.CreateUser(ctx.Request().Context(), *user)

	if err != nil {
		switch {
		case errors.Is(err, ErrUserExists):
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (uc *UserController) LoginUser(ctx *echo.Context) error {
	user, err := internal.BindAndValidate[LoginUserRequest](ctx)
	if err != nil {
		return err
	}

	resp, err := uc.service.LoginUser(ctx.Request().Context(), *user)
	if err != nil {
		// no matter the error, it is always safe to return an uniform error type
		return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidCredentials.Error())
	}

	issueTime := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(resp.ID)),
		IssuedAt:  jwt.NewNumericDate(issueTime),
		ExpiresAt: jwt.NewNumericDate(issueTime.Add(time.Hour * 8)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		// no matter the error, it is always safe to return an uniform error type
		return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidCredentials.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": signedToken,
	})
}

func (uc *UserController) GetUserInfo(ctx *echo.Context) error {
	return nil
}
