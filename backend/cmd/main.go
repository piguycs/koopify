package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"piguy.nl/koopify/internal"
	"piguy.nl/koopify/internal/db"
	"piguy.nl/koopify/internal/router"
	"piguy.nl/koopify/internal/user"
)

func main() {
	ctx := context.Background()
	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	validator := internal.NewCustomValidator()
	e.Validator = &validator

	e.Use(middleware.RequestLogger())
	// this API should be available for standalone use
	e.Use(middleware.CORS("*"))

	conn, err := pgx.Connect(ctx, config.PgDb)
	if err != nil {
		log.Fatal("unable to connect to postgres database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	userRepo := user.NewUserRepository(queries)
	userService := user.NewUserService(&userRepo)
	userController := user.NewUserController(config.JwtSecret, userService)

	router.RegisterPublicRoutes(e, &userController)
	router.RegisterPrivateRoutes(e, config.JwtSecret, &userController)

	for _, route := range e.Router().Routes() {
		log.Printf("Route: %s %s\n", route.Method, route.Path)
	}

	internal.StartServer(ctx, config.HostAddr, config.TlsConfig, e)
}
