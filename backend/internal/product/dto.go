package product

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"piguy.nl/koopify/internal/db"
)

// CategoryResponse is a lightweight category representation.
type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ProductResponse is the full product representation including its categories.
type ProductResponse struct {
	ID              int64              `json:"id"`
	Name            string             `json:"name"`
	Slug            string             `json:"slug"`
	Description     string             `json:"description"`
	ImageUrl        *string            `json:"imageUrl"`
	PriceEurCents   int32              `json:"priceEurCents"`
	DiscountPercent *int32             `json:"discountPercent"`
	InventoryCount  int32              `json:"inventoryCount"`
	InStock         bool               `json:"inStock"`
	IsActive        bool               `json:"isActive"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
	Categories      []CategoryResponse `json:"categories"`
}

// CreateProductRequest is the payload for creating a new product.
type CreateProductRequest struct {
	Name            string  `json:"name"            validate:"required,lte=255,gte=1"`
	Slug            string  `json:"slug"            validate:"required,lte=255,gte=1"`
	Description     string  `json:"description"     validate:"required"`
	ImageUrl        *string `json:"imageUrl"        validate:"omitempty,url"`
	PriceEurCents   int32   `json:"priceEurCents"   validate:"min=0"`
	DiscountPercent *int32  `json:"discountPercent" validate:"omitempty,min=0,max=100"`
	InventoryCount  int32   `json:"inventoryCount"  validate:"min=0"`
	InStock         bool    `json:"inStock"`
	IsActive        bool    `json:"isActive"`
	// Category IDs to attach; new category names can be created via the
	// category endpoints first and their IDs used here.
	CategoryIDs []int64 `json:"categoryIds"`
}

// UpdateProductRequest is the payload for updating an existing product.
type UpdateProductRequest struct {
	Name            string  `json:"name"            validate:"required,lte=255,gte=1"`
	Slug            string  `json:"slug"            validate:"required,lte=255,gte=1"`
	Description     string  `json:"description"     validate:"required"`
	ImageUrl        *string `json:"imageUrl"        validate:"omitempty,url"`
	PriceEurCents   int32   `json:"priceEurCents"   validate:"min=0"`
	DiscountPercent *int32  `json:"discountPercent" validate:"omitempty,min=0,max=100"`
	InventoryCount  int32   `json:"inventoryCount"  validate:"min=0"`
	InStock         bool    `json:"inStock"`
	IsActive        bool    `json:"isActive"`
	CategoryIDs     []int64 `json:"categoryIds"`
}

// CreateCategoryRequest is the payload for creating a new category.
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,lte=128,gte=1"`
	Slug string `json:"slug" validate:"required,lte=128,gte=1"`
}

func categoryResponseFrom(c db.Category) CategoryResponse {
	return CategoryResponse{
		ID:   c.ID,
		Name: c.Name,
		Slug: c.Slug,
	}
}

func productResponseFrom(p db.Product, categories []db.Category) ProductResponse {
	cats := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		cats[i] = categoryResponseFrom(c)
	}

	var imageUrl *string
	if p.ImageUrl.Valid {
		imageUrl = &p.ImageUrl.String
	}

	var discountPercent *int32
	if p.DiscountPercent.Valid {
		discountPercent = &p.DiscountPercent.Int32
	}

	return ProductResponse{
		ID:              p.ID,
		Name:            p.Name,
		Slug:            p.Slug,
		Description:     p.Description,
		ImageUrl:        imageUrl,
		PriceEurCents:   p.PriceEurCents,
		DiscountPercent: discountPercent,
		InventoryCount:  p.InventoryCount,
		InStock:         p.InStock,
		IsActive:        p.IsActive,
		CreatedAt:       p.CreatedAt.Time,
		UpdatedAt:       p.UpdatedAt.Time,
		Categories:      cats,
	}
}

func optionalText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func optionalInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: *i, Valid: true}
}
