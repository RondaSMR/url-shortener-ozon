package apperor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const name = "PS"

var (
	ErrInternalSystem = NewAppError(http.StatusInternalServerError, "00100", "internal system error")
	ErrBadRequest     = NewAppError(http.StatusBadRequest, "00101", "bad request error")
	ErrValidation     = NewAppError(http.StatusBadRequest, "00102", "validation error")
	ErrNotFound       = NewAppError(http.StatusNotFound, "00103", "resource not found")
	ErrUnauthorized   = NewAppError(http.StatusUnauthorized, "00104", "unauthorized")
	ErrForbidden      = NewAppError(http.StatusForbidden, "00105", "access forbidden")
	ErrInvalidID      = NewAppError(http.StatusBadRequest, "00106", "invalid id")
)

type ErrorFields map[string]string

type AppError struct {
	Err           error       `json:"-"`
	Message       string      `json:"message,omitempty"`
	Code          string      `json:"code,omitempty"`
	TransportCode int         `json:"-"`
	Fields        ErrorFields `json:"fields,omitempty"`
	TraceID       string      `json:"trace_id,omitempty"`
}

func (e *AppError) WithFields(fields ErrorFields) {
	e.Fields = fields
}

func NewAppError(transportCode int, code, message string) *AppError {
	return &AppError{
		Err:           fmt.Errorf(message),
		Code:          name + "-" + code,
		TransportCode: transportCode,
		Message:       message,
	}
}

func (e *AppError) Error() string {
	err := e.Err.Error()
	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			err += ", " + k + " " + v
		}
	}

	return err
}

func (e *AppError) JsonResponse(c *gin.Context, realError error) {
	logger := zap.L().WithOptions(zap.AddCallerSkip(1))
	logger.Error(e.Message, zap.Error(realError))
	c.JSON(e.TransportCode, e)
}

func (e *AppError) Unwrap() error { return e.Err }
