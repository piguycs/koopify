package checkout

import (
	"time"

	"piguy.nl/koopify/internal/db"
)

type CartItemRequest struct {
	ProductID int64 `json:"productId" validate:"required,min=1"`
	Quantity  int32 `json:"quantity" validate:"required,min=1,max=99"`
}

type CreateCheckoutSessionRequest struct {
	Items []CartItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderItemResponse struct {
	ID             int64  `json:"id"`
	ProductID      int64  `json:"productId"`
	ProductName    string `json:"productName"`
	Quantity       int32  `json:"quantity"`
	UnitPriceCents int32  `json:"unitPriceCents"`
}

type OrderResponse struct {
	ID                 int64               `json:"id"`
	UserID             int64               `json:"userId"`
	Status             string              `json:"status"`
	TotalEurCents      int32               `json:"totalEurCents"`
	AdyenPaymentLink   *string             `json:"adyenPaymentLink"`
	AdyenReference     *string             `json:"adyenReference"`
	AdyenSessionResult *string             `json:"adyenSessionResult"`
	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
	Items              []OrderItemResponse `json:"items"`
}

type UpdateAdyenSessionRequest struct {
	SessionId     string `json:"sessionId" validate:"required"`
	SessionResult string `json:"sessionResult" validate:"required"`
}

type CheckoutSessionResponse struct {
	Order            OrderResponse `json:"order"`
	AdyenCheckoutURL string        `json:"adyenCheckoutUrl"`
}

func orderItemResponseFrom(item db.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:             item.ID,
		ProductID:      item.ProductID,
		ProductName:    item.ProductName,
		Quantity:       item.Quantity,
		UnitPriceCents: item.UnitPriceCents,
	}
}

func orderResponseFrom(order db.Order, items []db.OrderItem) OrderResponse {
	itemResponses := make([]OrderItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = orderItemResponseFrom(item)
	}

	var adyenRef *string
	if order.AdyenReference.Valid {
		adyenRef = &order.AdyenReference.String
	}

	var adyenSessionResult *string
	if order.AdyenSessionResult.Valid {
		adyenSessionResult = &order.AdyenSessionResult.String
	}

	var adyenPaymentLink *string
	if order.AdyenPaymentLink.Valid {
		adyenPaymentLink = &order.AdyenPaymentLink.String
	}

	return OrderResponse{
		ID:                 order.ID,
		UserID:             order.UserID,
		Status:             order.Status,
		TotalEurCents:      order.TotalEurCents,
		AdyenPaymentLink:   adyenPaymentLink,
		AdyenReference:     adyenRef,
		AdyenSessionResult: adyenSessionResult,
		CreatedAt:          order.CreatedAt.Time,
		UpdatedAt:          order.UpdatedAt.Time,
		Items:              itemResponses,
	}
}
