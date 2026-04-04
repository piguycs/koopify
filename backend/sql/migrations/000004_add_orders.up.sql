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
