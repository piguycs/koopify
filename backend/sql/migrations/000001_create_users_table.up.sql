create table if not exists users (
    id           bigserial primary key,
    email        varchar   unique not null,
    display_name varchar   not null,
    password     varchar   not null,
    admin        boolean   not null default false
);
