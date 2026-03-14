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
