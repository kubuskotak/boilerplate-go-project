package usecase

import (
	"context"
	"database/sql"

	"github.com/kubuskotak/boilerplate-go-project/config"
	"github.com/kubuskotak/boilerplate-go-project/domain/entity"
	"github.com/kubuskotak/boilerplate-go-project/domain/repository"
	"github.com/kubuskotak/tyr"
	"github.com/kubuskotak/valkyrie"
	"github.com/rs/zerolog/log"
)

type RegisterParams struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func (s *Service) Register(ctx context.Context, args RegisterParams) (entity.User, error) {
	var user entity.User
	if err := s.store.WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		var txErr error
		txRepo := repository.NewSQL(tx)

		hashedPassword := valkyrie.HashPassword(args.Password, config.GetConfig().App.SecretKey)
		user, txErr = txRepo.CreateUser(ctx, repository.CreateUserParams{
			Username:       args.Username,
			FullName:       args.FullName,
			Email:          args.Email,
			HashedPassword: hashedPassword,
		})
		if txErr != nil {
			e := tyr.CatchErr(txErr)
			log.Error().Err(e).Msg("Register")
			return e
		}
		log.Info().Interface("user", user).Msg("Register: CreateUser")
		return nil
	}); err != nil {
		log.Error().Err(err).Msg("Register: WithTransaction")
		return user, err
	}
	return user, nil
}
