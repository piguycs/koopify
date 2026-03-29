package checkout

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"piguy.nl/koopify/internal/db"
)

type CheckoutRepository interface {
	GetProduct(ctx context.Context, id int64) (*db.Product, error)
	CreateOrder(ctx context.Context, params db.CreateOrderParams) (*db.Order, error)
	CreateOrderItem(ctx context.Context, params db.CreateOrderItemParams) (*db.OrderItem, error)
	GetOrder(ctx context.Context, id int64) (*db.Order, error)
	GetOrderByUser(ctx context.Context, orderID int64, userID int64) (*db.Order, error)
	ListOrderItems(ctx context.Context, orderID int64) ([]db.OrderItem, error)
	ListOrdersByUser(ctx context.Context, userID int64) ([]db.Order, error)
	ListAllOrders(ctx context.Context) ([]db.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) (*db.Order, error)
	UpdateOrderAdyenReference(ctx context.Context, orderID int64, adyenRef string) (*db.Order, error)
	UpdateOrderAdyenSession(ctx context.Context, orderID int64, sessionId string, sessionResult string) (*db.Order, error)
	DecrementProductInventory(ctx context.Context, productID int64, quantity int32) (*db.Product, error)
}

type PGCheckoutRepository struct {
	queries *db.Queries
}

func NewCheckoutRepository(queries *db.Queries) PGCheckoutRepository {
	return PGCheckoutRepository{queries: queries}
}

func (r PGCheckoutRepository) GetProduct(ctx context.Context, id int64) (*db.Product, error) {
	p, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r PGCheckoutRepository) CreateOrder(ctx context.Context, params db.CreateOrderParams) (*db.Order, error) {
	order, err := r.queries.CreateOrder(ctx, params)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r PGCheckoutRepository) CreateOrderItem(ctx context.Context, params db.CreateOrderItemParams) (*db.OrderItem, error) {
	item, err := r.queries.CreateOrderItem(ctx, params)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r PGCheckoutRepository) GetOrder(ctx context.Context, id int64) (*db.Order, error) {
	order, err := r.queries.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r PGCheckoutRepository) GetOrderByUser(ctx context.Context, orderID int64, userID int64) (*db.Order, error) {
	order, err := r.queries.GetOrderByUser(ctx, db.GetOrderByUserParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r PGCheckoutRepository) ListOrderItems(ctx context.Context, orderID int64) ([]db.OrderItem, error) {
	items, err := r.queries.ListOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r PGCheckoutRepository) ListOrdersByUser(ctx context.Context, userID int64) ([]db.Order, error) {
	orders, err := r.queries.ListOrdersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r PGCheckoutRepository) ListAllOrders(ctx context.Context) ([]db.Order, error) {
	orders, err := r.queries.ListAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r PGCheckoutRepository) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (*db.Order, error) {
	order, err := r.queries.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		ID:     orderID,
		Status: status,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r PGCheckoutRepository) UpdateOrderAdyenReference(ctx context.Context, orderID int64, adyenRef string) (*db.Order, error) {
	order, err := r.queries.UpdateOrderAdyenReference(ctx, db.UpdateOrderAdyenReferenceParams{
		ID:             orderID,
		AdyenReference: pgtype.Text{String: adyenRef, Valid: adyenRef != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r PGCheckoutRepository) UpdateOrderAdyenSession(ctx context.Context, orderID int64, sessionId string, sessionResult string) (*db.Order, error) {
	order, err := r.queries.UpdateOrderAdyenSession(ctx, db.UpdateOrderAdyenSessionParams{
		ID:                 orderID,
		AdyenReference:     pgtype.Text{String: sessionId, Valid: true},
		AdyenSessionResult: pgtype.Text{String: sessionResult, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r PGCheckoutRepository) DecrementProductInventory(ctx context.Context, productID int64, quantity int32) (*db.Product, error) {
	product, err := r.queries.DecrementProductInventory(ctx, db.DecrementProductInventoryParams{
		ID:             productID,
		InventoryCount: quantity,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInsufficientStock
		}
		return nil, err
	}
	return &product, nil
}
