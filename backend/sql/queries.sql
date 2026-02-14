-- name: GetUser :one
select * from users
    where id = $1 limit 1;

-- name: ListUsers :many
select * from users
    order by id;

-- name: CreateUser :one
insert into users (
    email, display_name, password
) values (
    $1, $2, $3
) returning *;

-- name: DeleteUser :exec
delete from users
    where id = $1;

-- name: GetUserWithEmail :one
select * from users
	where email = $1 limit 1;
