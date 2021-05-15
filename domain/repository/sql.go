package repository

import (
	"context"

	"github.com/kubuskotak/boilerplate-go-project/domain/entity"
	"github.com/kubuskotak/boilerplate-go-project/pkg/db"
)

type Adapter interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (entity.User, error)
	GetUser(ctx context.Context, username string) (entity.User, error)
}

type SQLStore struct {
	db db.Driver
}

func NewSQL(db db.Driver) Adapter {
	return &SQLStore{db: db}
}
