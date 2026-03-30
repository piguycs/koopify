package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"piguy.nl/koopify/internal"
	"piguy.nl/koopify/internal/checkout"
	"piguy.nl/koopify/internal/db"
	"piguy.nl/koopify/internal/product"
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

	charmLogHandler := log.New(os.Stderr)
	e.Logger = slog.New(charmLogHandler)
	e.Use(middleware.RequestLogger())
	// this API should be available for standalone use
	e.Use(middleware.CORS("*"))

	conn, err := pgxpool.New(ctx, config.PgDb)
	if err != nil {
		log.Fatal("unable to connect to postgres database", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	queries := db.New(conn)

	userRepo := user.NewUserRepository(queries)
	userService := user.NewUserService(&userRepo)
	userController := user.NewUserController(config.JwtSecret, userService)

	productRepo := product.NewProductRepository(queries)
	productService := product.NewProductService(&productRepo)
	productController := product.NewProductController(productService)

	checkoutRepo := checkout.NewCheckoutRepository(queries)
	checkoutService := checkout.NewCheckoutService(
		&checkoutRepo,
		config.AdyenApiKey,
		config.AdyenMerchantAccount,
		config.AdyenThemeId,
		config.CheckoutReturnUrl,
	)
	checkoutController := checkout.NewCheckoutController(checkoutService)

	router.RegisterPublicRoutes(e, &userController, &productController)
	router.RegisterPrivateRoutes(e, config.JwtSecret, &userController, &productController, &checkoutController)

	for _, route := range e.Router().Routes() {
		log.Info("Backend API route registered", "method", route.Method, "path", route.Path)
	}

	internal.StartServer(ctx, config.HostAddr, config.TlsConfig, e)
}
