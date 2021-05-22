package repository

import (
	"context"

	"github.com/kubuskotak/boilerplate-go-project/domain/entity"
	"github.com/kubuskotak/tyr"
)

type Adapter interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (entity.User, error)
	GetUser(ctx context.Context, username string) (entity.User, error)
}

type SQLStore struct {
	db tyr.Driver
}

func NewSQL(db tyr.Driver) Adapter {
	return &SQLStore{db: db}
}
