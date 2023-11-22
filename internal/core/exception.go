package core

import (
	"net/http"

	"github.com/goexl/exception"
	"github.com/goexl/gox"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Exception struct {
	_ gox.CannotCopy
}

func NewException() *Exception {
	return new(Exception)
}

func (e *Exception) Invalidation(err error) error {
	return e.error(http.StatusBadRequest, err)
}

func (e *Exception) Error(err error) error {
	return e.error(http.StatusInternalServerError, err)
}

func (e *Exception) New(code int, fields ...gox.Field[any]) (err error) {
	message := "服务器错误，客户端需要根据返回中的`code`码来确认具体是什么错误"
	exc := exception.New().Code(code).Message(message).Field(fields...).Build()
	err = e.error(http.StatusInternalServerError, exc)

	return
}

func (e *Exception) Notfound() error {
	return e.error(http.StatusNotFound, exception.New().Message("未找到资源").Build())
}

func (e *Exception) error(code codes.Code, err error) error {
	return status.Error(code, err.Error())
}
