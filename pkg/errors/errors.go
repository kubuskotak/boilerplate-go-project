package errors

import (
	"errors"
	"fmt"
)

const (
	StatusNoRows      = 1002
	SyntaxError       = 4601
	InvalidForeignKey = 4830
)

var (
	ErrSyntax            = errors.New("syntax error")
	ErrNoRows            = errors.New("no rows")
	ErrInvalidForeignKey = errors.New("invalid foreign key")
)

func Wrap(code int, err error, msg string) error {
	return &appError{
		code: code,
		err:  err,
		msg:  msg,
	}
}

type appError struct {
	code int
	err  error
	msg  string
}

func (e *appError) Unwrap() error { return e.err }

func (e *appError) Error() string {
	return fmt.Sprintf("%d: %v - %v", e.code, e.msg, e.err)
}

func InvalidForeignKeyErr(msg string) error {
	return Wrap(InvalidForeignKey, ErrInvalidForeignKey, msg)
}
