create table if not exists products (
    id               bigserial primary key,
    name             varchar   not null,
    slug             varchar   unique not null,
    description      text      not null,
    image_url        varchar,
    price_eur_cents  integer   not null check (price_eur_cents >= 0),
    discount_percent integer   check (discount_percent >= 0 and discount_percent <= 100),
    inventory_count  integer   not null default 0,
    in_stock         boolean   not null default false,
    is_active        boolean   not null default true,
    created_at       timestamptz not null default now(),
    updated_at       timestamptz not null default now()
);

create table if not exists categories (
    id          bigserial primary key,
    name        varchar   not null,
    slug        varchar   unique not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create table if not exists product_categories (
    product_id   bigint not null references products(id) on delete cascade,
    category_id  bigint not null references categories(id) on delete cascade,
    primary key (product_id, category_id)
);
