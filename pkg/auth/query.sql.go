// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getUserById = `-- name: GetUserById :one
select ID , logging, Nom,  Prenom, Roles from public.users where id = $1
`

type GetUserByIdRow struct {
	ID      int32
	Logging string
	Nom     pgtype.Text
	Prenom  pgtype.Text
	Roles   pgtype.Text
}

func (q *Queries) GetUserById(ctx context.Context, id int32) (GetUserByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.Logging,
		&i.Nom,
		&i.Prenom,
		&i.Roles,
	)
	return i, err
}

const getUserByLogging = `-- name: GetUserByLogging :one
select ID , logging, pwd_hash, Nom,  Prenom, Roles from public.users where logging = $1
`

type GetUserByLoggingRow struct {
	ID      int32
	Logging string
	PwdHash []byte
	Nom     pgtype.Text
	Prenom  pgtype.Text
	Roles   pgtype.Text
}

func (q *Queries) GetUserByLogging(ctx context.Context, logging string) (GetUserByLoggingRow, error) {
	row := q.db.QueryRow(ctx, getUserByLogging, logging)
	var i GetUserByLoggingRow
	err := row.Scan(
		&i.ID,
		&i.Logging,
		&i.PwdHash,
		&i.Nom,
		&i.Prenom,
		&i.Roles,
	)
	return i, err
}
