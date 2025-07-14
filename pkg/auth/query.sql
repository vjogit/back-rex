-- name: GetUserByLogging :one
select ID , logging, pwd_hash, Nom,  Prenom, Roles from public.users where logging = $1; 

-- name: GetUserById :one
select ID , logging, Nom,  Prenom, Roles from public.users where id = $1; 
