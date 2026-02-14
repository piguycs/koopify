package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"piguy.nl/koopify/internal"
	"piguy.nl/koopify/internal/db"
	"piguy.nl/koopify/internal/user"
)

var DefaultPgDb = "postgres://postgres:postgres@localhost:5432/?sslmode=disable"
var DefaultHostAddr = ":8080"
var CommitHash = "dev"

func main() {
	jwtSecret := internal.GetEnvDefault("JWT_SECRET", "INVALID_SECRET")
	if jwtSecret == "INVALID_SECRET" {
		panic("cannot run the application without a valid JWT secret")
	}

	ctx := context.Background()

	e := echo.New()
	app := e.Group("/api/v1", echojwt.JWT([]byte(jwtSecret)))
	appNoAuth := e.Group("/public_api/v1")

	validator := internal.NewCustomValidator()
	e.Validator = &validator

	e.Use(middleware.RequestLogger())
	// this API should be available for standalone use
	e.Use(middleware.CORS("*"))

	pgdb := internal.GetEnvDefault("PGDB", DefaultPgDb)
	conn, err := pgx.Connect(ctx, pgdb)
	if err != nil {
		log.Fatal("unable to connect to postgres database", "error", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	userRepo := user.NewUserRepository(queries)
	userService := user.NewUserService(&userRepo)
	userController := user.NewUserController(jwtSecret, userService)

	appNoAuth.POST("/login", userController.LoginUser)
	appNoAuth.POST("/sign_up", userController.CreateUser)
	app.GET("/users/me", userController.GetUserInfo)
	app.GET("/users/:id", userController.GetUserInfo)

	e.GET("/commit", func(c *echo.Context) error {
		return c.String(http.StatusOK, CommitHash)
	})

	// pass in the enviornment variable keys for reading the TLS config
	address := internal.GetEnvDefault("HOST_ADDR", DefaultHostAddr)
	tlsConfig, err := internal.TlsConfig("TLS_ENABLED", "TLS_CERT", "TLS_KEY")
	if err != nil {
		e.Logger.Error("failed to load TLS", "error", err)
	}

	internal.StartServer(ctx, address, tlsConfig, e)
}
