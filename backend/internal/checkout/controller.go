package checkout

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
	"piguy.nl/koopify/internal"
	"piguy.nl/koopify/internal/auth"
	"piguy.nl/koopify/internal/response"
)

type CheckoutController struct {
	service CheckoutService
}

func NewCheckoutController(service CheckoutService) CheckoutController {
	return CheckoutController{service: service}
}

func (cc *CheckoutController) CreateCheckoutSession(ctx *echo.Context) error {
	userID, err := auth.UserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.NewError("unauthorized", err.Error()))
	}

	req, err := internal.BindAndValidate[CreateCheckoutSessionRequest](ctx)
	if err != nil {
		return err
	}

	session, err := cc.service.CreateCheckoutSession(ctx.Request().Context(), userID, *req)
	if err != nil {
		switch {
		case errors.Is(err, ErrProductNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("product_not_found", err.Error()))
		case errors.Is(err, ErrInsufficientStock):
			return ctx.JSON(http.StatusBadRequest, response.NewError("insufficient_stock", err.Error()))
		default:
			if strings.Contains(err.Error(), "not available") {
				return ctx.JSON(http.StatusBadRequest, response.NewError("product_unavailable", err.Error()))
			}
			log.Errorf("Error creating checkout session: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to create checkout session"))
		}
	}

	return ctx.JSON(http.StatusCreated, session)
}

func (cc *CheckoutController) GetOrder(ctx *echo.Context) error {
	userID, err := auth.UserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.NewError("unauthorized", err.Error()))
	}

	orderID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.NewError("invalid_request", "invalid order id"))
	}

	order, err := cc.service.GetOrder(ctx.Request().Context(), userID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("order_not_found", err.Error()))
		default:
			log.Errorf("Error fetching order %d: %s", orderID, err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to fetch order"))
		}
	}

	return ctx.JSON(http.StatusOK, order)
}

func (cc *CheckoutController) ListOrders(ctx *echo.Context) error {
	userID, err := auth.UserIDFromToken(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.NewError("unauthorized", err.Error()))
	}

	orders, err := cc.service.ListOrders(ctx.Request().Context(), userID)
	if err != nil {
		log.Errorf("Error listing orders for user %d: %s", userID, err.Error())
		return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to list orders"))
	}

	return ctx.JSON(http.StatusOK, orders)
}

func (cc *CheckoutController) ListAllOrders(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	orders, err := cc.service.ListAllOrders(ctx.Request().Context())
	if err != nil {
		log.Errorf("Error listing all orders: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to list orders"))
	}

	return ctx.JSON(http.StatusOK, orders)
}
