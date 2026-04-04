create table if not exists password_resets (
    id               bigserial primary key,
    user_id          bigserial references users(id) not null,
    reset_token      varchar not null,
    created_at       timestamptz not null default now()
);
