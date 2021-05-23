package usecase

import (
	"context"
	"database/sql"
	"github.com/kubuskotak/tyr"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kubuskotak/boilerplate-go-project/config"
	"github.com/kubuskotak/boilerplate-go-project/domain/entity"
	"github.com/kubuskotak/boilerplate-go-project/domain/repository"

	"github.com/kubuskotak/valkyrie"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	service Service
	db      *tyr.Sql
	mock    sqlmock.Sqlmock
}

func TestUserMain(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (r *UserSuite) SetupTest() {
	data, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		require.Failf(r.T(), "failed to open stub db", "%v", err)
	}

	s := &tyr.Sql{
		Db: data,
	}
	r.service = Service{store: s}
	r.mock = mockDB
	r.db = s
	config.Instance = &config.Config{}
	if err := valkyrie.Config(valkyrie.ConfigOpts{
		Config:    config.Instance,
		Paths:     []string{"../../config"},
		Filenames: []string{"app.config.yaml", ".env"},
	}); err != nil {
		log.Error().Err(err).Msg("get config error")
	}
}

func (r *UserSuite) TearDownTest() {
	require.NoError(r.T(), r.mock.ExpectationsWereMet())
}

func (r *UserSuite) TestRegister() {
	t := r.T()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now().UTC()
	user := entity.User{
		Username:          "Suryakencana007",
		HashedPassword:    "pass1234",
		FullName:          "Surya kencana",
		Email:             "nanang.jobs@gmail.com",
		PasswordChangedAt: now,
		CreatedAt:         now,
	}
	args := RegisterParams{
		Username: user.Username,
		Password: user.HashedPassword,
		FullName: user.FullName,
		Email:    user.Email,
	}
	hashedPassword := valkyrie.HashPassword(args.Password, config.GetConfig().App.SecretKey)
	r.mock.ExpectBegin()
	r.mock.ExpectQuery(regexp.QuoteMeta(repository.CreateUserSql)).WithArgs(
		&args.Username,
		hashedPassword,
		&args.FullName,
		&args.Email,
	).WillReturnRows(
		sqlmock.NewRows([]string{
			"username", "hashed_password", "full_name",
			"email", "password_changed_at", "created_at"}).
			AddRow(user.Username, hashedPassword,
				user.FullName, user.Email,
				now, now))
	r.mock.ExpectCommit()
	u, err := r.service.Register(ctx, args)
	require.NoError(t, err)
	require.Equal(t, u.HashedPassword, hashedPassword)
}

func (r *UserSuite) TestRegisterTxFail() {
	t := r.T()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errTx := &tyr.Error{
		Code:    "1160",
		Message: sql.ErrNoRows.Error(),
	}

	now := time.Now().UTC()
	user := entity.User{
		Username:          "Suryakencana007",
		HashedPassword:    "pass1234",
		FullName:          "Surya kencana",
		Email:             "nanang.jobs@gmail.com",
		PasswordChangedAt: now,
		CreatedAt:         now,
	}
	args := RegisterParams{
		Username: user.Username,
		Password: user.HashedPassword,
		FullName: user.FullName,
		Email:    user.Email,
	}
	hashedPassword := valkyrie.HashPassword(args.Password, "sekret")
	r.mock.ExpectBegin()
	r.mock.ExpectQuery(regexp.QuoteMeta(repository.CreateUserSql)).WithArgs(
		&args.Username,
		hashedPassword,
		&args.FullName,
		&args.Email,
	).WillReturnError(errTx)
	r.mock.ExpectRollback()
	u, err := r.service.Register(ctx, args)
	require.Error(t, err)
	require.Equal(t, err.Error(), errTx.Error())
	require.Equal(t, u, entity.User{})
}
