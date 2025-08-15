package errs

import (
	"fmt"
)

type Error struct {
	Err  error
	Code int
	MSG  string
}

func (e *Error) Error() string {
	if e.MSG != "" {
		return fmt.Sprintf("%d; %s", e.Code, e.MSG)
	}

	return e.Err.Error()
}

func (e *Error) Msg() string {
	return e.MSG
}

func (e *Error) ErrCode() int {
	return e.Code
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Errf(err error, tmpl string, args ...interface{}) error {
	return &Error{Err: err, MSG: fmt.Sprintf(tmpl, args...)}
}

func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		MSG:  msg,
	}
}
