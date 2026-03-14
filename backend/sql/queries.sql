-- name: GetUser :one
select id, email, display_name, password, admin, requested_deletion_at, deletion_scheduled_at from users
    where id = $1 limit 1;

-- name: ListUsers :many
select id, email, display_name, password, admin, requested_deletion_at, deletion_scheduled_at from users
    order by id;

-- name: CreateUser :one
insert into users (
    email, display_name, password
) values (
    $1, $2, $3
) returning id, email, display_name, password, admin, requested_deletion_at, deletion_scheduled_at;

-- name: UpdateUser :one
update users
set display_name = $2,
    email = $3
where id = $1
returning id, email, display_name, password, admin, requested_deletion_at, deletion_scheduled_at;

-- name: DeleteUser :exec
delete from users
    where id = $1;

-- name: GetUserWithEmail :one
select id, email, display_name, password, admin, requested_deletion_at, deletion_scheduled_at from users
	where email = $1 limit 1;

-- name: GetDeletionPolicy :one
select id, deletion_delay_hours from deletion_policy
    limit 1;

-- name: RequestUserDeletion :one
update users
set requested_deletion_at = now(),
    deletion_scheduled_at = now() + make_interval(hours => $2)
where id = $1
returning id, email, display_name, password, admin, requested_deletion_at, deletion_scheduled_at;

-- name: CancelUserDeletion :one
update users
set requested_deletion_at = null,
    deletion_scheduled_at = null
where id = $1
returning id, email, display_name, password, admin, requested_deletion_at, deletion_scheduled_at;

-- Products
-- name: CreateProduct :one
insert into products (
    name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active
) values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) returning id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at;

-- name: GetProduct :one
select id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at
from products where id = $1 limit 1;

-- name: GetProductBySlug :one
select id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at
from products where slug = $1 limit 1;

-- name: ListProducts :many
select id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at
from products where is_active = true
order by created_at desc;

-- name: ListAllProducts :many
select id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at
from products order by created_at desc;

-- name: UpdateProduct :one
update products set
    name = $2,
    slug = $3,
    description = $4,
    image_url = $5,
    price_eur_cents = $6,
    discount_percent = $7,
    inventory_count = $8,
    in_stock = $9,
    is_active = $10,
    updated_at = now()
where id = $1
returning id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at;

-- name: DeleteProduct :exec
delete from products where id = $1;

-- name: ListProductsByCategory :many
select p.id, p.name, p.slug, p.description, p.image_url, p.price_eur_cents, p.discount_percent, p.inventory_count, p.in_stock, p.is_active, p.created_at, p.updated_at
from products p
join product_categories pc on p.id = pc.product_id
where pc.category_id = $1 and p.is_active = true
order by p.created_at desc;

-- Categories
-- name: CreateCategory :one
insert into categories (name, slug)
values ($1, $2)
returning id, name, slug, created_at, updated_at;

-- name: GetCategory :one
select id, name, slug, created_at, updated_at
from categories where id = $1 limit 1;

-- name: GetCategoryBySlug :one
select id, name, slug, created_at, updated_at
from categories where slug = $1 limit 1;

-- name: ListCategories :many
select id, name, slug, created_at, updated_at
from categories order by name;

-- name: UpdateCategory :one
update categories set
    name = $2,
    slug = $3,
    updated_at = now()
where id = $1
returning id, name, slug, created_at, updated_at;

-- name: DeleteCategory :exec
delete from categories where id = $1;

-- Product Category Links
-- name: AttachProductCategory :exec
insert into product_categories (product_id, category_id)
values ($1, $2) on conflict do nothing;

-- name: DetachProductCategory :exec
delete from product_categories where product_id = $1 and category_id = $2;

-- name: GetProductCategories :many
select c.id, c.name, c.slug, c.created_at, c.updated_at
from categories c
join product_categories pc on c.id = pc.category_id
where pc.product_id = $1
order by c.name;
