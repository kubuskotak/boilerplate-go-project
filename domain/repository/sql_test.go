package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kubuskotak/boilerplate-go-project/domain/entity"
	"github.com/kubuskotak/boilerplate-go-project/pkg/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepoSuite struct {
	suite.Suite
	store Adapter
	db    sqlmock.Sqlmock
}

func TestRepoMain(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}

func (r *RepoSuite) SetupTest() {
	data, mockDB, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		require.Failf(r.T(), "failed to open stub db", "%v", err)
	}

	r.store = NewSQL(data)
	r.db = mockDB
}

func (r *RepoSuite) TearDownTest() {
	require.NoError(r.T(), r.db.ExpectationsWereMet())
}

func (r *RepoSuite) TestCreateUser() {
	t := r.T()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now().UTC()
	user := entity.User{
		Username:          "Suryakencana007",
		HashedPassword:    "H34eslasjas78u8jknmalsikd",
		FullName:          "Surya kencana",
		Email:             "nanang.jobs@gmail.com",
		PasswordChangedAt: now,
		CreatedAt:         now,
	}

	actual := CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}
	r.db.ExpectQuery(regexp.QuoteMeta(CreateUserSql)).WithArgs(
		&actual.Username,
		&actual.HashedPassword,
		&actual.FullName,
		&actual.Email,
	).WillReturnRows(
		sqlmock.NewRows([]string{
			"username", "hashed_password", "full_name",
			"email", "password_changed_at", "created_at"}).
			AddRow(actual.Username, actual.HashedPassword,
				actual.FullName, actual.Email,
				now, now))
	expected, err := r.store.CreateUser(ctx, actual)
	require.NoError(t, err)
	require.Equal(t, expected, user)
}

func (r *RepoSuite) TestCreateUserFail() {
	t := r.T()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now().UTC()
	user := entity.User{
		Username:          "Suryakencana007",
		HashedPassword:    "H34eslasjas78u8jknmalsikd",
		FullName:          "Surya kencana",
		Email:             "nanang.jobs@gmail.com",
		PasswordChangedAt: now,
		CreatedAt:         now,
	}

	actual := &db.Error{
		Code:    "1160",
		Message: sql.ErrNoRows.Error(),
	}

	params := CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}
	r.db.ExpectQuery(regexp.QuoteMeta(CreateUserSql)).WithArgs(
		&params.Username,
		&params.HashedPassword,
		&params.FullName,
		&params.Email,
	).WillReturnError(actual)

	expected, err := r.store.CreateUser(ctx, params)
	require.Error(t, err)
	require.Equal(t, err.Error(), actual.Error())
	require.Equal(t, expected, entity.User{})
}

func (r *RepoSuite) TestGetUser() {
	t := r.T()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now().UTC()
	user := entity.User{
		Username:          "Suryakencana007",
		HashedPassword:    "H34eslasjas78u8jknmalsikd",
		FullName:          "Surya kencana",
		Email:             "nanang.jobs@gmail.com",
		PasswordChangedAt: now,
		CreatedAt:         now,
	}

	actual := CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}
	r.db.ExpectQuery(regexp.QuoteMeta(getUser)).WithArgs(
		&actual.Username,
	).WillReturnRows(
		sqlmock.NewRows([]string{
			"username", "hashed_password", "full_name",
			"email", "password_changed_at", "created_at"}).
			AddRow(actual.Username, actual.HashedPassword,
				actual.FullName, actual.Email,
				now, now))
	expected, err := r.store.GetUser(ctx, actual.Username)
	require.NoError(t, err)
	require.Equal(t, expected, user)
}

func (r *RepoSuite) TestGetUserFail() {
	t := r.T()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now().UTC()
	user := entity.User{
		Username:          "Suryakencana007",
		HashedPassword:    "H34eslasjas78u8jknmalsikd",
		FullName:          "Surya kencana",
		Email:             "nanang.jobs@gmail.com",
		PasswordChangedAt: now,
		CreatedAt:         now,
	}

	actual := &db.Error{
		Code:    "1160",
		Message: sql.ErrNoRows.Error(),
	}

	params := CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}
	r.db.ExpectQuery(regexp.QuoteMeta(getUser)).WithArgs(
		&params.Username,
	).WillReturnError(actual)

	expected, err := r.store.GetUser(ctx, params.Username)
	require.Error(t, err)
	require.Equal(t, err.Error(), actual.Error())
	require.Equal(t, expected, entity.User{})
}
