package router

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"piguy.nl/koopify/internal/auth"
	"piguy.nl/koopify/internal/user"
)

func RegisterPublicRoutes(e *echo.Echo, userController *user.UserController) {
	public := e.Group("/public_api/v1")
	public.POST("/login", userController.LoginUser)
	public.POST("/sign_up", userController.CreateUser)
	public.GET("/password_policy", userController.GetPasswordPolicy)
}

func RegisterPrivateRoutes(e *echo.Echo, jwtSecret string, userController *user.UserController) {
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret),
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(auth.KoopifyClaims)
		},
	})

	private := e.Group("/api/v1", jwtMiddleware)
	private.GET("/users/me", userController.GetCurrentUser)
	private.PATCH("/users/me", userController.UpdateCurrentUser)
	private.POST("/users/me/deletion", userController.RequestDeletion)
	private.DELETE("/users/me/deletion", userController.CancelDeletion)
	private.GET("/users", userController.ListUsers)
	private.GET("/users/:id", userController.GetUserByID)
	private.PATCH("/users/:id/admin", userController.UpdateUserAdmin)
	private.POST("/users/:id/reset_password", userController.TriggerPasswordReset)
	private.POST("/users/:id/deletion", userController.RequestUserDeletionAdmin)
}
