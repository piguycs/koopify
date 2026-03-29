create table users (
    id           bigserial primary key,
    email        varchar   unique not null,
    display_name varchar   not null,
    password     varchar   not null,
    admin        boolean   not null default false,
    requested_deletion_at timestamptz,
    deletion_scheduled_at timestamptz
);

create table deletion_policy (
    id boolean primary key default true,
    deletion_delay_hours integer not null default 24,
    check (id)
);

create table products (
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

create table categories (
    id          bigserial primary key,
    name        varchar   not null,
    slug        varchar   unique not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now()
);

create table product_categories (
    product_id   bigint not null references products(id) on delete cascade,
    category_id  bigint not null references categories(id) on delete cascade,
    primary key (product_id, category_id)
);

create table orders (
    id              bigserial primary key,
    user_id         bigint not null references users(id),
    status          varchar not null default 'pending',
    total_eur_cents integer not null check (total_eur_cents >= 0),
    adyen_reference varchar,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now()
);

create table order_items (
    id              bigserial primary key,
    order_id        bigint not null references orders(id) on delete cascade,
    product_id      bigint not null references products(id),
    product_name    varchar not null,
    quantity        integer not null check (quantity > 0),
    unit_price_cents integer not null check (unit_price_cents >= 0),
    created_at      timestamptz not null default now()
);
