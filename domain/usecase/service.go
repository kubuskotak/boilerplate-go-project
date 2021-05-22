package usecase

import (
	"github.com/kubuskotak/tyr"
)

type Service struct {
	store *tyr.Sql
}

func NewService(sql *tyr.Sql) *Service {
	return &Service{store: sql}
}
