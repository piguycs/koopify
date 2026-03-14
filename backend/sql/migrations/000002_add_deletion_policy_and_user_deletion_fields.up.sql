-- boolean ID with check means we only store a singular row in this table
create table deletion_policy (
    id boolean primary key default true,
    deletion_delay_hours integer not null default 24,
    check (id)
);

insert into deletion_policy (id) values (true);

alter table users add column requested_deletion_at timestamptz;
alter table users add column deletion_scheduled_at timestamptz;
