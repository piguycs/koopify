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

-- name: UpdateUserPassword :one
update users
set password = $2
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

-- name: UpdateUserAdmin :one
update users
set admin = $2
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
select 
	id,
	name,
	slug,
	description,
	image_url,
	price_eur_cents,
	discount_percent,
	inventory_count,
	in_stock,
	is_active,
	created_at,
	updated_at
from products where slug = $1 limit 1;

-- name: ListProducts :many
select id,
	name,
	slug,
	description,
	image_url,
	price_eur_cents,
	discount_percent,
	inventory_count,
	in_stock,
	is_active,
	created_at,
	updated_at
from products where is_active = true
order by created_at desc;

-- name: ListAllProducts :many
select
	id,
	name,
	slug,
	description,
	image_url,
	price_eur_cents,
	discount_percent,
	inventory_count,
	in_stock,
	is_active,
	created_at,
	updated_at
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
select
	p.id,
	p.name,
	p.slug,
	p.description,
	p.image_url,
	p.price_eur_cents,
	p.discount_percent,
	p.inventory_count,
	p.in_stock,
	p.is_active,
	p.created_at,
	p.updated_at
from products p
join product_categories pc on p.id = pc.product_id
where pc.category_id = $1 and p.is_active = true
order by p.created_at desc;

-- name: ListProductsPaginated :many
select id,
	name,
	slug,
	description,
	image_url,
	price_eur_cents,
	discount_percent,
	inventory_count,
	in_stock,
	is_active,
	created_at,
	updated_at
from products
where is_active = true
  and (coalesce($3, '') = '' or name ilike '%' || $3 || '%' or description ilike '%' || $3 || '%')
order by created_at desc
limit $1 offset $2;

-- name: ListProductsPaginatedByCategory :many
select
	p.id,
	p.name,
	p.slug,
	p.description,
	p.image_url,
	p.price_eur_cents,
	p.discount_percent,
	p.inventory_count,
	p.in_stock,
	p.is_active,
	p.created_at,
	p.updated_at
from products p
join product_categories pc on p.id = pc.product_id
where pc.category_id = $1
  and p.is_active = true
  and (coalesce($4, '') = '' or p.name ilike '%' || $4 || '%' or p.description ilike '%' || $4 || '%')
order by p.created_at desc
limit $2 offset $3;

-- name: CountActiveProducts :one
select count(*) from products
where is_active = true
  and (coalesce($1, '') = '' or name ilike '%' || $1 || '%' or description ilike '%' || $1 || '%');

-- name: CountActiveProductsByCategory :one
select count(*)
from products p
join product_categories pc on p.id = pc.product_id
where pc.category_id = $1
  and p.is_active = true
  and (coalesce($2, '') = '' or p.name ilike '%' || $2 || '%' or p.description ilike '%' || $2 || '%');

-- Admin: List all products (including inactive) with pagination and search
-- name: ListAllProductsPaginated :many
select id,
	name,
	slug,
	description,
	image_url,
	price_eur_cents,
	discount_percent,
	inventory_count,
	in_stock,
	is_active,
	created_at,
	updated_at
from products
where (coalesce($3, '') = '' or name ilike '%' || $3 || '%' or description ilike '%' || $3 || '%')
order by created_at desc
limit $1 offset $2;

-- Admin: List all products by category with pagination and search
-- name: ListAllProductsPaginatedByCategory :many
select
	p.id,
	p.name,
	p.slug,
	p.description,
	p.image_url,
	p.price_eur_cents,
	p.discount_percent,
	p.inventory_count,
	p.in_stock,
	p.is_active,
	p.created_at,
	p.updated_at
from products p
join product_categories pc on p.id = pc.product_id
where pc.category_id = $1
  and (coalesce($4, '') = '' or p.name ilike '%' || $4 || '%' or p.description ilike '%' || $4 || '%')
order by p.created_at desc
limit $2 offset $3;

-- Admin: Count all products (including inactive) with search
-- name: CountAllProducts :one
select count(*) from products
where (coalesce($1, '') = '' or name ilike '%' || $1 || '%' or description ilike '%' || $1 || '%');

-- Admin: Count all products by category with search
-- name: CountAllProductsByCategory :one
select count(*)
from products p
join product_categories pc on p.id = pc.product_id
where pc.category_id = $1
  and (coalesce($2, '') = '' or p.name ilike '%' || $2 || '%' or p.description ilike '%' || $2 || '%');

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

-- Orders
-- name: CreateOrder :one
insert into orders (
	user_id,
	status,
	total_eur_cents,
	adyen_payment_link,
	adyen_reference,
	adyen_session_result
) values ($1, $2, $3, $4, $5, $6)
returning id,
		  user_id,
		  status,
		  total_eur_cents,
		  adyen_payment_link,
		  adyen_reference,
		  adyen_session_result,
		  created_at,
		  updated_at;

-- name: GetOrder :one
select id,
	   user_id,
	   status,
	   total_eur_cents,
	   adyen_payment_link,
	   adyen_reference,
	   adyen_session_result,
	   created_at,
	   updated_at
from orders where id = $1 limit 1;

-- name: GetOrderByUser :one
select id, user_id, status, total_eur_cents, adyen_payment_link, adyen_reference, adyen_session_result, created_at, updated_at
from orders where id = $1 and user_id = $2 limit 1;

-- name: ListOrdersByUser :many
select id, user_id, status, total_eur_cents, adyen_payment_link, adyen_reference, adyen_session_result, created_at, updated_at
from orders where user_id = $1 order by created_at desc;

-- Admin: List all orders
-- name: ListAllOrders :many
select id, user_id, status, total_eur_cents, adyen_payment_link, adyen_reference, adyen_session_result, created_at, updated_at
from orders order by created_at desc;

-- name: UpdateOrderStatus :one
update orders set status = $2, updated_at = now()
where id = $1
returning id, user_id, status, total_eur_cents, adyen_payment_link, adyen_reference, adyen_session_result, created_at, updated_at;

-- name: UpdateOrderPaymentLink :one
update orders set adyen_payment_link = $2, updated_at = now()
where id = $1
returning id, user_id, status, total_eur_cents, adyen_payment_link, adyen_reference, adyen_session_result, created_at, updated_at;

-- name: UpdateOrderAdyenReference :one
update orders set adyen_reference = $2, updated_at = now()
where id = $1
returning id, user_id, status, total_eur_cents, adyen_payment_link, adyen_reference, adyen_session_result, created_at, updated_at;

-- name: UpdateOrderAdyenSession :one
update orders set adyen_reference = $2, adyen_session_result = $3, updated_at = now()
where id = $1
returning id, user_id, status, total_eur_cents, adyen_payment_link, adyen_reference, adyen_session_result, created_at, updated_at;

-- name: CreateOrderItem :one
insert into order_items (order_id, product_id, product_name, quantity, unit_price_cents)
values ($1, $2, $3, $4, $5)
returning id, order_id, product_id, product_name, quantity, unit_price_cents, created_at;

-- name: ListOrderItems :many
select id, order_id, product_id, product_name, quantity, unit_price_cents, created_at
from order_items where order_id = $1 order by created_at;

-- name: GetOrderItem :one
select id, order_id, product_id, product_name, quantity, unit_price_cents, created_at
from order_items where id = $1 limit 1;

-- Inventory management
-- name: DecrementProductInventory :one
update products set
    inventory_count = inventory_count - $2,
    in_stock = (inventory_count - $2) > 0,
    updated_at = now()
where id = $1 and inventory_count >= $2
returning id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at;

-- name: GetProductForUpdate :one
select id, name, slug, description, image_url, price_eur_cents, discount_percent, inventory_count, in_stock, is_active, created_at, updated_at
from products where id = $1 limit 1 for update;
