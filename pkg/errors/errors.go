package errors

import (
	stderrors "errors"
	"fmt"
)

type Code string

const (
	CodeInvalidArgument     Code = "INVALID_ARGUMENT"
	CodeAlreadyExists       Code = "ALREADY_EXISTS"
	CodeNotFound            Code = "NOT_FOUND"
	CodeConflict            Code = "CONFLICT"
	CodePartialFailed       Code = "PARTIAL_FAILED"
	CodeVersionConflict     Code = "VERSION_CONFLICT"
	CodeWorkspaceTombstoned Code = "WORKSPACE_TOMBSTONED"
	CodeWorkspaceConflicted Code = "WORKSPACE_CONFLICTED"
	CodeProviderUnsupported Code = "PROVIDER_UNSUPPORTED"
	CodeProviderUnavailable Code = "PROVIDER_UNAVAILABLE"
	CodeValidationFailed    Code = "VALIDATION_FAILED"
	CodeQueryParseError     Code = "QUERY_PARSE_ERROR"
	CodeQueryPlanError      Code = "QUERY_PLAN_ERROR"
	CodeToolDisabled        Code = "TOOL_DISABLED"
	CodeToolNotFound        Code = "TOOL_NOT_FOUND"
	CodeTimeout             Code = "TIMEOUT"
	CodeInternal            Code = "INTERNAL"
	CodeNotImplemented      Code = "NOT_IMPLEMENTED"
)

type Error struct {
	Code      Code              `json:"code"`
	Message   string            `json:"message"`
	Retryable bool              `json:"retryable"`
	Details   map[string]string `json:"details,omitempty"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return string(e.Code)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func WithDetails(code Code, message string, details map[string]string) *Error {
	return &Error{Code: code, Message: message, Details: details}
}

func Retryable(code Code, message string) *Error {
	return &Error{Code: code, Message: message, Retryable: true}
}

func As(err error) (*Error, bool) {
	var target *Error
	if stderrors.As(err, &target) {
		return target, true
	}
	return nil, false
}

func IsCode(err error, code Code) bool {
	if target, ok := As(err); ok {
		return target.Code == code
	}
	return false
}
