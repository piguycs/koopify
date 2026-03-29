package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"piguy.nl/koopify/internal/db"
)

// ProductRepository maps one-to-one with SQL queries. No composition or
// business logic lives here; the service layer owns that.
type ProductRepository interface {
	ListProducts(ctx context.Context) ([]db.Product, error)
	ListAllProducts(ctx context.Context) ([]db.Product, error)
	ListProductsByCategory(ctx context.Context, categoryID int64) ([]db.Product, error)
	GetProduct(ctx context.Context, id int64) (*db.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*db.Product, error)
	CreateProduct(ctx context.Context, params db.CreateProductParams) (*db.Product, error)
	UpdateProduct(ctx context.Context, params db.UpdateProductParams) (*db.Product, error)
	DeleteProduct(ctx context.Context, id int64) error

	GetProductCategories(ctx context.Context, productID int64) ([]db.Category, error)
	AttachProductCategory(ctx context.Context, params db.AttachProductCategoryParams) error
	DetachProductCategory(ctx context.Context, params db.DetachProductCategoryParams) error

	ListCategories(ctx context.Context) ([]db.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*db.Category, error)
	CreateCategory(ctx context.Context, params db.CreateCategoryParams) (*db.Category, error)

	// Pagination methods
	ListProductsPaginated(ctx context.Context, limit int32, offset int32, searchTerm string) ([]db.Product, error)
	ListProductsPaginatedByCategory(ctx context.Context, categoryID int64, limit int32, offset int32, searchTerm string) ([]db.Product, error)
	CountActiveProducts(ctx context.Context, searchTerm string) (int64, error)
	CountActiveProductsByCategory(ctx context.Context, categoryID int64, searchTerm string) (int64, error)
}

// PGProductRepository is the PostgreSQL implementation of ProductRepository.
type PGProductRepository struct {
	queries *db.Queries
}

func NewProductRepository(queries *db.Queries) PGProductRepository {
	return PGProductRepository{queries: queries}
}

func (r PGProductRepository) ListProducts(ctx context.Context) ([]db.Product, error) {
	return r.queries.ListProducts(ctx)
}

func (r PGProductRepository) ListAllProducts(ctx context.Context) ([]db.Product, error) {
	return r.queries.ListAllProducts(ctx)
}

func (r PGProductRepository) ListProductsByCategory(
	ctx context.Context,
	categoryID int64,
) ([]db.Product, error) {
	return r.queries.ListProductsByCategory(ctx, categoryID)
}

func (r PGProductRepository) GetProduct(ctx context.Context, id int64) (*db.Product, error) {
	p, err := r.queries.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r PGProductRepository) GetProductBySlug(ctx context.Context, slug string) (*db.Product, error) {
	p, err := r.queries.GetProductBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r PGProductRepository) CreateProduct(
	ctx context.Context,
	params db.CreateProductParams,
) (*db.Product, error) {
	p, err := r.queries.CreateProduct(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	return &p, nil
}

func (r PGProductRepository) UpdateProduct(
	ctx context.Context,
	params db.UpdateProductParams,
) (*db.Product, error) {
	p, err := r.queries.UpdateProduct(ctx, params)
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
	return &p, nil
}

func (r PGProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	return r.queries.DeleteProduct(ctx, id)
}

func (r PGProductRepository) GetProductCategories(
	ctx context.Context,
	productID int64,
) ([]db.Category, error) {
	return r.queries.GetProductCategories(ctx, productID)
}

func (r PGProductRepository) AttachProductCategory(
	ctx context.Context,
	params db.AttachProductCategoryParams,
) error {
	return r.queries.AttachProductCategory(ctx, params)
}

func (r PGProductRepository) DetachProductCategory(
	ctx context.Context,
	params db.DetachProductCategoryParams,
) error {
	return r.queries.DetachProductCategory(ctx, params)
}

func (r PGProductRepository) ListCategories(ctx context.Context) ([]db.Category, error) {
	return r.queries.ListCategories(ctx)
}

func (r PGProductRepository) GetCategoryBySlug(ctx context.Context, slug string) (*db.Category, error) {
	c, err := r.queries.GetCategoryBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r PGProductRepository) CreateCategory(
	ctx context.Context,
	params db.CreateCategoryParams,
) (*db.Category, error) {
	c, err := r.queries.CreateCategory(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrSlugCategoryTaken
		}
		return nil, err
	}
	return &c, nil
}

func (r PGProductRepository) ListProductsPaginated(
	ctx context.Context,
	limit int32,
	offset int32,
	searchTerm string,
) ([]db.Product, error) {
	return r.queries.ListProductsPaginated(ctx, db.ListProductsPaginatedParams{
		Limit:   limit,
		Offset:  offset,
		Column3: searchTerm,
	})
}

func (r PGProductRepository) ListProductsPaginatedByCategory(
	ctx context.Context,
	categoryID int64,
	limit int32,
	offset int32,
	searchTerm string,
) ([]db.Product, error) {
	return r.queries.ListProductsPaginatedByCategory(ctx, db.ListProductsPaginatedByCategoryParams{
		CategoryID: categoryID,
		Limit:      limit,
		Offset:     offset,
		Column4:    searchTerm,
	})
}

func (r PGProductRepository) CountActiveProducts(ctx context.Context, searchTerm string) (int64, error) {
	return r.queries.CountActiveProducts(ctx, searchTerm)
}

func (r PGProductRepository) CountActiveProductsByCategory(ctx context.Context, categoryID int64, searchTerm string) (int64, error) {
	return r.queries.CountActiveProductsByCategory(ctx, db.CountActiveProductsByCategoryParams{
		CategoryID: categoryID,
		Column2:    searchTerm,
	})
}
