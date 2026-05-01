package error

const (
	CodeSuccess             = 0
	CodeInternalServerError = 1000
	CodeBadRequest          = 1001
	CodeUnauthorized        = 1002
	CodeTooManyRequests     = 1003
	CodeUnimplemented       = 1004
	CodeForbidden           = 1005
	CodeDataNotFound        = 1006
)

var (
	NonError             = NewError(CodeSuccess, "请求成功")
	InternalServerError  = NewError(CodeInternalServerError, "服务器出错, 请稍后重试")
	InvalidParamsError   = NewError(CodeBadRequest, "参数有误")
	UnauthorizedError    = NewError(CodeUnauthorized, "未授权登录")
	TooManyRequestsError = NewError(CodeTooManyRequests, "请求过于频繁")
	UnimplementedError   = NewError(CodeUnimplemented, "该操作尚未实现")
	ForbiddenError       = NewError(CodeForbidden, "禁止该操作")
	DataNotFoundError    = NewError(CodeDataNotFound, "数据不存在")
)
