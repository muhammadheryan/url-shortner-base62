package errors

import "github.com/muhammadheryan/url-shortner-base62/constant"

type CustomError struct {
	errType constant.ErrorType
}

func (c CustomError) Error() string {
	return constant.ErrorTypeMessage[c.errType]
}

func (c CustomError) ErrorCode() string {
	return constant.ErrorTypeCode[c.errType]
}

func (c CustomError) ErrorHTTPCode() int {
	return constant.ErrorTypeHTTPCode[c.errType]
}

func SetCustomError(errorType constant.ErrorType) CustomError {
	return CustomError{
		errType: errorType,
	}
}
