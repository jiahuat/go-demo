package wscperror

var (
	ErrOne = NewErr(1, "err one", nil)
	ErrTwo = NewErr(2, "err two", nil)
)

type wscpErr struct {
	Code int
	Msg  string
	Err  error
}

func NewErr(code int, msg string, err error) error {
	return &wscpErr{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
}

func (e *wscpErr) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return e.Msg
	}

	return e.Msg + e.Err.Error()
}

func (e *wscpErr) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
