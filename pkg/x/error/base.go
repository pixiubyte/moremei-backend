package error

type Error struct {
	code  int
	msg   string
	cause string
}

func (e *Error) GetCode() int {
	return e.code
}

func (e *Error) GetMessage() string {
	return e.msg
}

func (e *Error) GetCause() string {
	return e.cause
}

func (e *Error) Error() string {
	if len(e.cause) > 0 {
		return e.msg + "：" + e.cause
	}
	return e.msg
}

func NewError(code int, msg string) *Error {
	return &Error{code: code, msg: msg}
}

func WithCause(e *Error, cause string) *Error {
	e.cause = cause
	return e
}
