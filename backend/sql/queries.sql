-- name: GetUser :one
select id, email, display_name, password, admin from users
    where id = $1 limit 1;

-- name: ListUsers :many
select id, email, display_name, password, admin from users
    order by id;

-- name: CreateUser :one
insert into users (
    email, display_name, password
) values (
    $1, $2, $3
) returning id, email, display_name, password, admin;

-- name: DeleteUser :exec
delete from users
    where id = $1;

-- name: GetUserWithEmail :one
select id, email, display_name, password, admin from users
	where email = $1 limit 1;
