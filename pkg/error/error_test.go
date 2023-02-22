package wscperror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWscperror(t *testing.T) {
	var cases = []struct {
		err    error
		target error
		want   bool
	}{
		{err: nil, target: nil, want: true},
		{err: NewErr(1, "msg", nil), target: ErrOne, want: false},
		{err: NewErr(1, "msg", ErrOne), target: ErrOne, want: true},
		{err: NewErr(1, "msg", ErrTwo), target: ErrOne, want: false},
		{err: NewErr(1, "err one", nil), target: ErrOne, want: false},
		{err: ErrOne, target: ErrOne, want: true},
	}
	for _, c := range cases {
		got := errors.Is(c.err, c.target)
		assert.Equal(t, c.want, got)
	}
	// as true
	var wErr *wscpErr
	rErr := NewErr(3, "msg", nil)
	as := errors.As(rErr, &wErr)
	assert.Equal(t, true, as)
	assert.Equal(t, 3, wErr.Code)
	assert.Equal(t, "msg", wErr.Msg)

	asT := errors.As(ErrOne, &wErr)
	assert.Equal(t, true, asT)

	// as false
	asF := errors.As(errors.New("some err"), &wErr)
	assert.Equal(t, false, asF)

}
