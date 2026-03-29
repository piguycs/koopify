package router

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"piguy.nl/koopify/internal/auth"
	"piguy.nl/koopify/internal/checkout"
	"piguy.nl/koopify/internal/product"
	"piguy.nl/koopify/internal/user"
)

func RegisterPublicRoutes(
	e *echo.Echo,
	userController *user.UserController,
	productController *product.ProductController,
) {
	public := e.Group("/public_api/v1")
	public.POST("/login", userController.LoginUser)
	public.POST("/sign_up", userController.CreateUser)
	public.GET("/password_policy", userController.GetPasswordPolicy)

	public.GET("/products", productController.ListProducts)
	public.GET("/products/:slug", productController.GetProductBySlug)
	public.GET("/categories", productController.ListCategories)
}

func RegisterPrivateRoutes(
	e *echo.Echo,
	jwtSecret string,
	userController *user.UserController,
	productController *product.ProductController,
	checkoutController *checkout.CheckoutController,
) {
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
	private.DELETE("/users/:id/deletion", userController.CancelUserDeletionAdmin)
	private.PATCH("/users/:id", userController.UpdateUserDetailsAdmin)

	// Admin product routes (auth enforced per-handler)
	// Note: public GET /products and GET /products/:slug live in the public router
	private.GET("/products", productController.ListAllProductsPaginated)
	private.GET("/products/:id", productController.GetProduct)
	private.POST("/products", productController.CreateProduct)
	private.PUT("/products/:id", productController.UpdateProduct)
	private.DELETE("/products/:id", productController.DeleteProduct)

	// only admin users can access this. the GET endpoint is in the public router fn
	private.POST("/categories", productController.CreateCategory)

	// Checkout routes
	private.POST("/checkout/sessions", checkoutController.CreateCheckoutSession)
	private.GET("/orders", checkoutController.ListOrders)
	private.GET("/orders/:id", checkoutController.GetOrder)

	// Admin order routes (admin check enforced in handler)
	private.GET("/admin/orders", checkoutController.ListAllOrders)
}
