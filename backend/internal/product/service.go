package product

import "context"

// ProductService defines the business-logic layer for products and categories.
type ProductService interface {
	ListAllProducts(ctx context.Context) ([]ProductResponse, error)
	GetProduct(ctx context.Context, id int64) (*ProductResponse, error)
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, id int64, req UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, id int64) error

	ListCategories(ctx context.Context) ([]CategoryResponse, error)
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error)
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) ListAllProducts(ctx context.Context) ([]ProductResponse, error) {
	return s.repo.ListAllProducts(ctx)
}

func (s *productService) GetProduct(ctx context.Context, id int64) (*ProductResponse, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *productService) CreateProduct(
	ctx context.Context,
	req CreateProductRequest,
) (*ProductResponse, error) {
	return s.repo.CreateProduct(ctx, req)
}

func (s *productService) UpdateProduct(
	ctx context.Context,
	id int64,
	req UpdateProductRequest,
) (*ProductResponse, error) {
	return s.repo.UpdateProduct(ctx, id, req)
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *productService) ListCategories(ctx context.Context) ([]CategoryResponse, error) {
	return s.repo.ListCategories(ctx)
}

func (s *productService) CreateCategory(
	ctx context.Context,
	req CreateCategoryRequest,
) (*CategoryResponse, error) {
	return s.repo.CreateCategory(ctx, req)
}
