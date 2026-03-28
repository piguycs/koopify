package product

import "errors"

const pgUniqueViolation = "23505"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrSlugTaken         = errors.New("a product with this slug already exists")
	ErrCategoryNotFound  = errors.New("category not found")
	ErrSlugCategoryTaken = errors.New("a category with this slug already exists")
)
