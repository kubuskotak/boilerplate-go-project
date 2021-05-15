package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

type SqlStore interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}

type Driver interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
}

const (
	POSTGRES string = "postgres"
	MYSQL    string = "mysql"
)

type Sql struct {
	Db *sql.DB
	SqlStore
}

func (r *Sql) WithTransaction(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		if errRoll := tx.Rollback(); errRoll != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, errRoll)
		}
		return err
	}
	return tx.Commit()
}

type SqlParams struct {
	Driver      string
	DSN         string
	MaxOpen     int
	MaxIdle     int
	MaxLifeTime int
}

type SqlConnParams struct {
	Driver, Dsn string
}

func New(args SqlConnParams) (*Sql, error) {
	db, err := sql.Open(args.Driver, args.Dsn)
	if err != nil {
		panic(fmt.Errorf("cannot access your db master connection").Error())
	}

	return &Sql{Db: db}, nil
}

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error %s: %s", e.Code, e.Message)
}

func CatchErr(err error) *Error {
	var e *Error
	// specific database error
	// postgre
	if pqErr, ok := err.(*pq.Error); ok {
		e = &Error{
			Code:    pqErr.Code.Name(),
			Message: pqErr.Message,
		}
	}
	// mysql
	if myErr, ok := err.(*mysql.MySQLError); ok {
		e = &Error{
			Code:    strconv.Itoa(int(myErr.Number)),
			Message: myErr.Message,
		}
	}
	// default
	if er, ok := err.(*Error); ok {
		e = er
	}

	return e
}
