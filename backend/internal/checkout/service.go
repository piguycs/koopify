package checkout

import (
	"context"
	"fmt"

	"github.com/adyen/adyen-go-api-library/v21/src/adyen"
	adyen_checkout "github.com/adyen/adyen-go-api-library/v21/src/checkout"
	"github.com/adyen/adyen-go-api-library/v21/src/common"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"piguy.nl/koopify/internal/db"
)

type CheckoutService struct {
	repo            CheckoutRepository
	apiKey          string
	merchantAccount string
	themeId         string
	returnUrl       string
	adyenClient     *adyen.APIClient
}

func NewCheckoutService(
	repo CheckoutRepository,
	apiKey, merchantAccount, themeId, returnUrl string,
) CheckoutService {
	cs := CheckoutService{
		repo:            repo,
		apiKey:          apiKey,
		merchantAccount: merchantAccount,
		themeId:         themeId,
		returnUrl:       returnUrl,
	}

	client := adyen.NewClient(&common.Config{
		ApiKey:      cs.apiKey,
		Environment: common.TestEnv,
	})

	cs.adyenClient = client

	return cs
}

type CartItemForCheckout struct {
	ProductID   int64
	ProductName string
	Quantity    int32
	UnitPrice   int32
}

func (s *CheckoutService) CreateCheckoutSession(
	ctx context.Context,
	userID int64,
	req CreateCheckoutSessionRequest,
) (*CheckoutSessionResponse, error) {
	cartItems := make([]CartItemForCheckout, 0, len(req.Items))
	totalPrice := int32(0)

	for _, item := range req.Items {
		product, err := s.repo.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d: %w", item.ProductID, err)
		}

		if !product.IsActive {
			return nil, fmt.Errorf("product %d: not available", item.ProductID)
		}

		if product.InventoryCount < item.Quantity {
			return nil, fmt.Errorf("product %d: insufficient stock (available: %d, requested: %d)",
				item.ProductID, product.InventoryCount, item.Quantity)
		}

		unitPrice := product.PriceEurCents
		if product.DiscountPercent.Valid {
			unitPrice = int32(float64(unitPrice) * (1 - float64(product.DiscountPercent.Int32)/100))
		}

		cartItems = append(cartItems, CartItemForCheckout{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
		})

		totalPrice += unitPrice * item.Quantity
	}

	order, err := s.repo.CreateOrder(ctx, db.CreateOrderParams{
		UserID:         userID,
		Status:         "pending",
		TotalEurCents:  totalPrice,
		AdyenReference: pgtype.Text{Valid: false},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range cartItems {
		_, err = s.repo.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:        order.ID,
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPrice,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	checkoutURL, err := s.createAdyenSession(ctx, order.ID, totalPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to create adyen session: %w", err)
	}

	items, err := s.repo.ListOrderItems(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	return &CheckoutSessionResponse{
		Order:            orderResponseFrom(*order, items),
		AdyenCheckoutURL: checkoutURL,
	}, nil
}

func (s *CheckoutService) createAdyenSession(ctx context.Context, orderID int64, amount int32) (string, error) {
	service := s.adyenClient.Checkout()

	amountEur := adyen_checkout.Amount{
		Currency: "EUR",
		Value:    int64(amount),
	}

	createCheckoutSessionRequest := adyen_checkout.CreateCheckoutSessionRequest{
		Reference:       fmt.Sprintf("order_%d", orderID),
		Mode:            common.PtrString("hosted"),
		Amount:          amountEur,
		MerchantAccount: s.merchantAccount,
		CountryCode:     common.PtrString("NL"),
		ThemeId:         common.PtrString(s.themeId),
		ReturnUrl:       s.returnUrl,
	}

	req := service.PaymentsApi.
		SessionsInput().
		IdempotencyKey(uuid.New().String()).
		CreateCheckoutSessionRequest(createCheckoutSessionRequest)

	res, httpRes, err := service.PaymentsApi.Sessions(ctx, req)
	_ = httpRes

	if err != nil {
		log.Error("Could not complete adyen request", "error", err)
		return "", ErrCheckoutCreationFailed
	}

	if res.Url == nil || *res.Url == "" {
		log.Error("adyen's response url was empty")
		return "", ErrCheckoutCreationFailed
	}

	return *res.Url, nil
}

func (s *CheckoutService) GetOrder(ctx context.Context, userID int64, orderID int64) (*OrderResponse, error) {
	order, err := s.repo.GetOrderByUser(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.ListOrderItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	resp := orderResponseFrom(*order, items)
	return &resp, nil
}

func (s *CheckoutService) ListOrders(ctx context.Context, userID int64) ([]OrderResponse, error) {
	orders, err := s.repo.ListOrdersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		items, err := s.repo.ListOrderItems(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch order items for order %d: %w", order.ID, err)
		}
		responses[i] = orderResponseFrom(order, items)
	}

	return responses, nil
}

func (s *CheckoutService) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (*OrderResponse, error) {
	order, err := s.repo.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.ListOrderItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}

	resp := orderResponseFrom(*order, items)
	return &resp, nil
}
