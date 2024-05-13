package service

import (
	"net/http"
)

const (
	code500 = http.StatusInternalServerError
	code400 = http.StatusBadRequest
	code401 = http.StatusUnauthorized

	InternalError       = "internal error"
	UserNotExist        = "user doesn't exist"
	AuthMessageNotExist = "auth message doesn't exist"
	ParseTokenFailed    = "parse token failed"

	TokenWrongSecret   = "wrong token secret"
	AuthMessageExpired = "auth message expired"
	WrongSignature     = "wrong signature"
	EcrecoverFailed    = "ecrecover failed"
	CreateUserFailed   = "create user failed"
)

// error struct
type ServiceError struct {
	Err    error  `json:"-"`
	Code   int64  `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
	Msg    string `json:"message,omitempty"`
}

func (e *ServiceError) Error() string {

	return e.Err.Error()
}

func (e *ServiceError) Unwrap() error { return e.Err }

func newServiceError(code int, err error, msg, detail string) *ServiceError {
	return &ServiceError{
		Code:   int64(code),
		Err:    err,
		Msg:    msg,
		Detail: detail,
	}
}
