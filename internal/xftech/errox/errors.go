package errox

type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

func NewBizError(code int, msg string) *BizError {
	return &BizError{Code: code, Message: msg}
}
