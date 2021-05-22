package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrNoRow_errorsIs(t *testing.T) {
	errNoRows := Wrap(StatusNoRows, ErrNoRows, "sql: 1601")
	require.ErrorIs(t, errNoRows, ErrNoRows)
	errSyntax := Wrap(SyntaxError, ErrSyntax, "sql: 1601")
	require.ErrorIs(t, errSyntax, ErrSyntax)
	errInvalidForeignKey := InvalidForeignKeyErr("sql: column is ")
	require.ErrorIs(t, errInvalidForeignKey, ErrInvalidForeignKey)
	t.Log(errInvalidForeignKey)
}
