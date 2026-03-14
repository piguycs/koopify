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
