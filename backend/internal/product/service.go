package product

import (
	"context"

	"piguy.nl/koopify/internal/db"
)

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return ProductService{repo: repo}
}

func (s *ProductService) ListActiveProducts(ctx context.Context) ([]ProductResponse, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	return s.enrichProducts(ctx, products)
}

func (s *ProductService) ListActiveProductsByCategorySlug(
	ctx context.Context,
	categorySlug string,
) ([]ProductResponse, error) {
	cat, err := s.repo.GetCategoryBySlug(ctx, categorySlug)
	if err != nil {
		return nil, err
	}
	products, err := s.repo.ListProductsByCategory(ctx, cat.ID)
	if err != nil {
		return nil, err
	}
	return s.enrichProducts(ctx, products)
}

func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*ProductResponse, error) {
	p, err := s.repo.GetProductBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return s.enrichProduct(ctx, p)
}

func (s *ProductService) ListAllProducts(ctx context.Context) ([]ProductResponse, error) {
	products, err := s.repo.ListAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	return s.enrichProducts(ctx, products)
}

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*ProductResponse, error) {
	p, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.enrichProduct(ctx, p)
}

func (s *ProductService) CreateProduct(
	ctx context.Context,
	req CreateProductRequest,
) (*ProductResponse, error) {
	p, err := s.repo.CreateProduct(ctx, db.CreateProductParams{
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
		return nil, err
	}

	if err = s.setProductCategories(ctx, p.ID, req.CategoryIDs); err != nil {
		return nil, err
	}

	return s.enrichProduct(ctx, p)
}

func (s *ProductService) UpdateProduct(
	ctx context.Context,
	id int64,
	req UpdateProductRequest,
) (*ProductResponse, error) {
	// Verify product exists before updating.
	_, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	p, err := s.repo.UpdateProduct(ctx, db.UpdateProductParams{
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
		return nil, err
	}

	if err = s.setProductCategories(ctx, p.ID, req.CategoryIDs); err != nil {
		return nil, err
	}

	return s.enrichProduct(ctx, p)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	// Verify product exists before deleting.
	_, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) ListCategories(ctx context.Context) ([]CategoryResponse, error) {
	cats, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]CategoryResponse, len(cats))
	for i, c := range cats {
		result[i] = categoryResponseFrom(c)
	}
	return result, nil
}

func (s *ProductService) CreateCategory(
	ctx context.Context,
	req CreateCategoryRequest,
) (*CategoryResponse, error) {
	c, err := s.repo.CreateCategory(ctx, db.CreateCategoryParams{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		return nil, err
	}
	resp := categoryResponseFrom(*c)
	return &resp, nil
}

func (s *ProductService) ListActiveProductsPaginated(
	ctx context.Context,
	start int32,
	end int32,
	searchTerm string,
) (*ProductResponsePage, error) {
	if start < 0 {
		start = 0
	}
	if end <= start {
		end = start
	}

	limit := end - start
	offset := start

	products, err := s.repo.ListProductsPaginated(ctx, limit, offset, searchTerm)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.repo.CountActiveProducts(ctx, searchTerm)
	if err != nil {
		return nil, err
	}

	enriched, err := s.enrichProducts(ctx, products)
	if err != nil {
		return nil, err
	}

	return &ProductResponsePage{
		Start:         start,
		End:           int32(len(enriched)),
		TotalProducts: totalProducts,
		Products:      enriched,
	}, nil
}

// ListActiveProductsPaginatedByCategory returns active products in a category with index-based pagination.
func (s *ProductService) ListActiveProductsPaginatedByCategory(
	ctx context.Context,
	categorySlug string,
	start int32,
	end int32,
	searchTerm string,
) (*ProductResponsePage, error) {
	cat, err := s.repo.GetCategoryBySlug(ctx, categorySlug)
	if err != nil {
		return nil, err
	}

	if start < 0 {
		start = 0
	}
	if end <= start {
		end = start
	}

	limit := end - start
	offset := start

	products, err := s.repo.ListProductsPaginatedByCategory(ctx, cat.ID, limit, offset, searchTerm)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.repo.CountActiveProductsByCategory(ctx, cat.ID, searchTerm)
	if err != nil {
		return nil, err
	}

	enriched, err := s.enrichProducts(ctx, products)
	if err != nil {
		return nil, err
	}

	return &ProductResponsePage{
		Start:         start,
		End:           int32(len(enriched)),
		TotalProducts: totalProducts,
		Products:      enriched,
	}, nil
}

// Admin: ListAllProductsPaginated returns all products (including inactive) with pagination and search.
func (s *ProductService) ListAllProductsPaginated(
	ctx context.Context,
	start int32,
	end int32,
	searchTerm string,
) (*ProductResponsePage, error) {
	if start < 0 {
		start = 0
	}
	if end <= start {
		end = start
	}

	limit := end - start
	offset := start

	products, err := s.repo.ListAllProductsPaginated(ctx, limit, offset, searchTerm)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.repo.CountAllProducts(ctx, searchTerm)
	if err != nil {
		return nil, err
	}

	enriched, err := s.enrichProducts(ctx, products)
	if err != nil {
		return nil, err
	}

	return &ProductResponsePage{
		Start:         start,
		End:           int32(len(enriched)),
		TotalProducts: totalProducts,
		Products:      enriched,
	}, nil
}

// Admin: ListAllProductsPaginatedByCategory returns all products in a category with pagination and search.
func (s *ProductService) ListAllProductsPaginatedByCategory(
	ctx context.Context,
	categorySlug string,
	start int32,
	end int32,
	searchTerm string,
) (*ProductResponsePage, error) {
	cat, err := s.repo.GetCategoryBySlug(ctx, categorySlug)
	if err != nil {
		return nil, err
	}

	if start < 0 {
		start = 0
	}
	if end <= start {
		end = start
	}

	limit := end - start
	offset := start

	products, err := s.repo.ListAllProductsPaginatedByCategory(ctx, cat.ID, limit, offset, searchTerm)
	if err != nil {
		return nil, err
	}

	totalProducts, err := s.repo.CountAllProductsByCategory(ctx, cat.ID, searchTerm)
	if err != nil {
		return nil, err
	}

	enriched, err := s.enrichProducts(ctx, products)
	if err != nil {
		return nil, err
	}

	return &ProductResponsePage{
		Start:         start,
		End:           int32(len(enriched)),
		TotalProducts: totalProducts,
		Products:      enriched,
	}, nil
}

// replaces the full category set for a producs. It removes categories which are not wanted and
// adds new ones
func (s *ProductService) setProductCategories(
	ctx context.Context,
	productID int64,
	categoryIDs []int64,
) error {
	existing, err := s.repo.GetProductCategories(ctx, productID)
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

	for _, c := range existing {
		if _, keep := desiredSet[c.ID]; !keep {
			err = s.repo.DetachProductCategory(ctx, db.DetachProductCategoryParams{
				ProductID:  productID,
				CategoryID: c.ID,
			})
			if err != nil {
				return err
			}
		}
	}

	for _, id := range categoryIDs {
		if _, already := existingSet[id]; !already {
			err = s.repo.AttachProductCategory(ctx, db.AttachProductCategoryParams{
				ProductID:  productID,
				CategoryID: id,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// fetches the categories for a product and builds a ProductResponse.
func (s *ProductService) enrichProduct(ctx context.Context, p *db.Product) (*ProductResponse, error) {
	cats, err := s.repo.GetProductCategories(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	resp := productResponseFrom(*p, cats)
	return &resp, nil
}

// fetches categories for each product and builds a `[]ProductResponse`
func (s *ProductService) enrichProducts(ctx context.Context, products []db.Product) ([]ProductResponse, error) {
	result := make([]ProductResponse, len(products))
	for i := range products {
		enriched, err := s.enrichProduct(ctx, &products[i])
		if err != nil {
			return nil, err
		}
		result[i] = *enriched
	}
	return result, nil
}
