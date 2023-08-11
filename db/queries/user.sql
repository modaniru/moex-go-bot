-- name: CreateUser :exec
insert into users (id) values ($1);

-- name: GetUser :one
select * from users where id = $1 limit 1;

-- name: DeleteUser :exec
delete from users where id = $1;