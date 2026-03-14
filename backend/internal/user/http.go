package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"piguy.nl/koopify/internal"
	"piguy.nl/koopify/internal/auth"
	"piguy.nl/koopify/internal/response"
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
			return ctx.JSON(http.StatusConflict, response.NewError("user_exists", err.Error()))
		default:
			log.Errorf("Error when trying to create user: %s\n", err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to create user"))
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
		return ctx.JSON(http.StatusUnauthorized, response.NewError("invalid_credentials", ErrInvalidCredentials.Error()))
	}

	issueTime := time.Now()
	claims := auth.KoopifyClaims{
		Admin: resp.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(resp.ID, 10),
			IssuedAt:  jwt.NewNumericDate(issueTime),
			ExpiresAt: jwt.NewNumericDate(issueTime.Add(time.Hour * 8)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		// no matter the error, it is always safe to return an uniform error type
		return ctx.JSON(http.StatusUnauthorized, response.NewError("invalid_credentials", ErrInvalidCredentials.Error()))
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": signedToken,
	})
}

func (uc *UserController) GetCurrentUser(ctx *echo.Context) error {
	userID, err := auth.UserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.NewError("unauthorized", "invalid auth token"))
	}

	resp, err := uc.service.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("user_not_found", err.Error()))
		default:
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to fetch user"))
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (uc *UserController) GetUserByID(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.NewError("invalid_request", "invalid user id"))
	}

	resp, err := uc.service.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("user_not_found", err.Error()))
		default:
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to fetch user"))
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}
