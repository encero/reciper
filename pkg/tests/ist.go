package tests

import (
	"testing"

	"github.com/matryer/is"
)

func newIst(t *testing.T) IsT {
	return IsT{
		Is: is.New(t),
		T:  t,
	}
}

type IsT struct {
	Is *is.I
	T  *testing.T
}

func (is IsT) NoErr(err error) {
	is.Is.Helper()
	is.Is.NoErr(err)
}

func (is IsT) Equal(a any, b any) {
	is.Is.Helper()
	is.Is.Equal(a, b)
}

func (is IsT) True(a bool) {
	is.Is.Helper()
	is.Is.True(a)
}
