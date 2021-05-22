package repository

import (
	"context"

	"github.com/kubuskotak/boilerplate-go-project/domain/entity"
	"github.com/kubuskotak/tyr"
	"github.com/rs/zerolog/log"
)

const CreateUserSql = `-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING username, hashed_password, full_name, email, password_changed_at, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (q *SQLStore) CreateUser(ctx context.Context, arg CreateUserParams) (entity.User, error) {
	row := q.db.QueryRowContext(ctx, CreateUserSql,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var u entity.User
	err := row.Scan(
		&u.Username,
		&u.HashedPassword,
		&u.FullName,
		&u.Email,
		&u.PasswordChangedAt,
		&u.CreatedAt,
	)
	if err != nil {
		e := tyr.CatchErr(err)
		log.Error().Err(e).Msg("CreateUser")
		return u, e
	}
	return u, err
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, password_changed_at, created_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *SQLStore) GetUser(ctx context.Context, username string) (entity.User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var u entity.User
	err := row.Scan(
		&u.Username,
		&u.HashedPassword,
		&u.FullName,
		&u.Email,
		&u.PasswordChangedAt,
		&u.CreatedAt,
	)
	if err != nil {
		e := tyr.CatchErr(err)
		log.Error().Err(e).Msg("GetUser")
		return u, e
	}
	return u, nil
}
