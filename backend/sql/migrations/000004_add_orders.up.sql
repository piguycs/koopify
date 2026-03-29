CREATE TABLE orders (
    id              bigserial PRIMARY KEY,
    user_id         bigint NOT NULL REFERENCES users(id),
    status          varchar NOT NULL DEFAULT 'pending',
    total_eur_cents integer NOT NULL CHECK (total_eur_cents >= 0),
    adyen_reference varchar,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
    id              bigserial PRIMARY KEY,
    order_id        bigint NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id      bigint NOT NULL REFERENCES products(id),
    product_name    varchar NOT NULL,
    quantity        integer NOT NULL CHECK (quantity > 0),
    unit_price_cents integer NOT NULL CHECK (unit_price_cents >= 0),
    created_at      timestamptz NOT NULL DEFAULT now()
);
