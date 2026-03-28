package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"
	"piguy.nl/koopify/internal"
	"piguy.nl/koopify/internal/auth"
	"piguy.nl/koopify/internal/response"
)

type ProductController struct {
	service ProductService
}

func NewProductController(service ProductService) ProductController {
	return ProductController{service: service}
}

// ListProducts is a public handler. If a ?category=<slug> query param is provided it returns only
// active products in that category (404 if the category slug is unknown). Without the param it
// returns all active products.
func (pc *ProductController) ListProducts(ctx *echo.Context) error {
	categorySlug := ctx.QueryParam("category")

	if categorySlug != "" {
		products, err := pc.service.ListActiveProductsByCategorySlug(ctx.Request().Context(), categorySlug)
		if err != nil {
			switch {
			case errors.Is(err, ErrCategoryNotFound):
				return ctx.JSON(http.StatusNotFound, response.NewError("category_not_found", err.Error()))
			default:
				log.Errorf("Error listing products by category %q: %s", categorySlug, err.Error())
				return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to list products"))
			}
		}
		return ctx.JSON(http.StatusOK, products)
	}

	products, err := pc.service.ListActiveProducts(ctx.Request().Context())
	if err != nil {
		log.Errorf("Error listing products: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to list products"))
	}
	return ctx.JSON(http.StatusOK, products)
}

// Returns a product by its slug, including inactive ones. So bookmarks etc dont appear as dead
// links to customers
func (pc *ProductController) GetProductBySlug(ctx *echo.Context) error {
	slug := ctx.Param("slug")

	product, err := pc.service.GetProductBySlug(ctx.Request().Context(), slug)
	if err != nil {
		switch {
		case errors.Is(err, ErrProductNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("product_not_found", err.Error()))
		default:
			log.Errorf("Error fetching product by slug %q: %s", slug, err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to fetch product"))
		}
	}
	return ctx.JSON(http.StatusOK, product)
}

// ListAllProducts is an admin handler that returns all products even if it is inactive. If someone
// has a direct link to this, maybe via a bookmark, they should not be left wondering what it was for
func (pc *ProductController) ListAllProducts(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	products, err := pc.service.ListAllProducts(ctx.Request().Context())
	if err != nil {
		log.Errorf("Error listing all products: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to list products"))
	}
	return ctx.JSON(http.StatusOK, products)
}

// GetProduct is an admin handler that returns a single product by numeric ID.
func (pc *ProductController) GetProduct(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.NewError("invalid_request", "invalid product id"))
	}

	product, err := pc.service.GetProduct(ctx.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrProductNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("product_not_found", err.Error()))
		default:
			log.Errorf("Error fetching product %d: %s", id, err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to fetch product"))
		}
	}
	return ctx.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product.
func (pc *ProductController) CreateProduct(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	req, err := internal.BindAndValidate[CreateProductRequest](ctx)
	if err != nil {
		return err
	}

	product, err := pc.service.CreateProduct(ctx.Request().Context(), *req)
	if err != nil {
		switch {
		case errors.Is(err, ErrSlugTaken):
			return ctx.JSON(http.StatusConflict, response.NewError("slug_taken", err.Error()))
		default:
			log.Errorf("Error creating product: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to create product"))
		}
	}

	return ctx.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product by ID.
func (pc *ProductController) UpdateProduct(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.NewError("invalid_request", "invalid product id"))
	}

	req, err := internal.BindAndValidate[UpdateProductRequest](ctx)
	if err != nil {
		return err
	}

	product, err := pc.service.UpdateProduct(ctx.Request().Context(), id, *req)
	if err != nil {
		switch {
		case errors.Is(err, ErrProductNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("product_not_found", err.Error()))
		case errors.Is(err, ErrSlugTaken):
			return ctx.JSON(http.StatusConflict, response.NewError("slug_taken", err.Error()))
		default:
			log.Errorf("Error updating product %d: %s", id, err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to update product"))
		}
	}

	return ctx.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID.
func (pc *ProductController) DeleteProduct(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.NewError("invalid_request", "invalid product id"))
	}

	err = pc.service.DeleteProduct(ctx.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrProductNotFound):
			return ctx.JSON(http.StatusNotFound, response.NewError("product_not_found", err.Error()))
		default:
			log.Errorf("Error deleting product %d: %s", id, err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to delete product"))
		}
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

// ListCategories is a public handler that returns all categories.
func (pc *ProductController) ListCategories(ctx *echo.Context) error {
	cats, err := pc.service.ListCategories(ctx.Request().Context())
	if err != nil {
		log.Errorf("Error listing categories: %s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to list categories"))
	}

	return ctx.JSON(http.StatusOK, cats)
}

// CreateCategory creates a new category.
func (pc *ProductController) CreateCategory(ctx *echo.Context) error {
	if !auth.IsAdminFromToken(ctx) {
		return ctx.JSON(http.StatusForbidden, response.NewError("forbidden", "admin access required"))
	}

	req, err := internal.BindAndValidate[CreateCategoryRequest](ctx)
	if err != nil {
		return err
	}

	cat, err := pc.service.CreateCategory(ctx.Request().Context(), *req)
	if err != nil {
		switch {
		case errors.Is(err, ErrSlugCategoryTaken):
			return ctx.JSON(http.StatusConflict, response.NewError("slug_taken", err.Error()))
		default:
			log.Errorf("Error creating category: %s", err.Error())
			return ctx.JSON(http.StatusInternalServerError, response.NewError("internal_error", "failed to create category"))
		}
	}

	return ctx.JSON(http.StatusCreated, cat)
}
