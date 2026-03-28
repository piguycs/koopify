package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"piguy.nl/koopify/internal/db"
)

// ProductRepository defines all persistence operations for products and categories.
type ProductRepository interface {
	ListAllProducts(ctx context.Context) ([]ProductResponse, error)
	GetProduct(ctx context.Context, id int64) (*ProductResponse, error)
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id int64) error

	ListCategories(ctx context.Context) ([]CategoryResponse, error)
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error)

	SetProductCategories(ctx context.Context, productID int64, categoryIDs []int64) error
}

// PGProductRepository is the PostgreSQL implementation of ProductRepository.
type PGProductRepository struct {
	queries *db.Queries
}

func NewProductRepository(queries *db.Queries) PGProductRepository {
	return PGProductRepository{queries: queries}
}

func (r PGProductRepository) ListAllProducts(ctx context.Context) ([]ProductResponse, error) {
	products, err := r.queries.ListAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]ProductResponse, len(products))
	for i, p := range products {
		cats, err := r.queries.GetProductCategories(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		result[i] = productResponseFrom(p, cats)
	}
	return result, nil
}

func (r PGProductRepository) GetProduct(ctx context.Context, id int64) (*ProductResponse, error) {
	p, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	cats, err := r.queries.GetProductCategories(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	resp := productResponseFrom(p, cats)
	return &resp, nil
}

func (r PGProductRepository) CreateProduct(
	ctx context.Context,
	req CreateProductRequest,
) (*ProductResponse, error) {
	p, err := r.queries.CreateProduct(ctx, db.CreateProductParams{
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		ImageUrl:        optionalText(req.ImageUrl),
		PriceEurCents:   req.PriceEurCents,
		DiscountPercent: optionalInt4(req.DiscountPercent),
		InventoryCount:  req.InventoryCount,
		InStock:         req.InStock,
		IsActive:        req.IsActive,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrSlugTaken
		}
		return nil, err
	}

	if err = r.SetProductCategories(ctx, p.ID, req.CategoryIDs); err != nil {
		return nil, err
	}

	cats, err := r.queries.GetProductCategories(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	resp := productResponseFrom(p, cats)
	return &resp, nil
}

func (r PGProductRepository) UpdateProduct(
	ctx context.Context,
	id int64,
	req UpdateProductRequest,
) (*ProductResponse, error) {
	p, err := r.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:              id,
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		ImageUrl:        optionalText(req.ImageUrl),
		PriceEurCents:   req.PriceEurCents,
		DiscountPercent: optionalInt4(req.DiscountPercent),
		InventoryCount:  req.InventoryCount,
		InStock:         req.InStock,
		IsActive:        req.IsActive,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrSlugTaken
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if err = r.SetProductCategories(ctx, p.ID, req.CategoryIDs); err != nil {
		return nil, err
	}

	cats, err := r.queries.GetProductCategories(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	resp := productResponseFrom(p, cats)
	return &resp, nil
}

func (r PGProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	// Verify the product exists first.
	_, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}
	return r.queries.DeleteProduct(ctx, id)
}

func (r PGProductRepository) ListCategories(ctx context.Context) ([]CategoryResponse, error) {
	cats, err := r.queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]CategoryResponse, len(cats))
	for i, c := range cats {
		result[i] = categoryResponseFrom(c)
	}
	return result, nil
}

func (r PGProductRepository) CreateCategory(
	ctx context.Context,
	req CreateCategoryRequest,
) (*CategoryResponse, error) {
	c, err := r.queries.CreateCategory(ctx, db.CreateCategoryParams{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrSlugCategoryTaken
		}
		return nil, err
	}
	resp := categoryResponseFrom(c)
	return &resp, nil
}

// SetProductCategories replaces the full set of categories for a product.
// It detaches all existing links and re-attaches the provided IDs.
func (r PGProductRepository) SetProductCategories(
	ctx context.Context,
	productID int64,
	categoryIDs []int64,
) error {
	existing, err := r.queries.GetProductCategories(ctx, productID)
	if err != nil {
		return err
	}

	existingSet := make(map[int64]struct{}, len(existing))
	for _, c := range existing {
		existingSet[c.ID] = struct{}{}
	}

	desiredSet := make(map[int64]struct{}, len(categoryIDs))
	for _, id := range categoryIDs {
		desiredSet[id] = struct{}{}
	}

	// Detach categories that are no longer desired.
	for _, c := range existing {
		if _, keep := desiredSet[c.ID]; !keep {
			if err = r.queries.DetachProductCategory(ctx, db.DetachProductCategoryParams{
				ProductID:  productID,
				CategoryID: c.ID,
			}); err != nil {
				return err
			}
		}
	}

	// Attach new categories.
	for _, id := range categoryIDs {
		if _, already := existingSet[id]; !already {
			if err = r.queries.AttachProductCategory(ctx, db.AttachProductCategoryParams{
				ProductID:  productID,
				CategoryID: id,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
