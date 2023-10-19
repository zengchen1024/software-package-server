package allerror

import "strings"

const (
	errorCodeNoPermission = "no_permission"

	ErrorCodeAccessTokenMissing = "access_token_missing"
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	ErrorCodeTooManyRequest   = "too_many_request"
	ErrorCodeSensitiveContent = "sensitive_content"

	ErrorCodeNotImporter    = "software_pkg_not_importer"
	ErrorCodeCIIsRunning    = "software_pkg_ci_is_running"
	ErrorCodeIncorrectPhase = "software_pkg_incorrect_phase"
)

// errorImpl
type errorImpl struct {
	code string
	msg  string
}

func (e errorImpl) Error() string {
	return e.msg
}

func (e errorImpl) ErrorCode() string {
	return e.code
}

// New
func New(code string, msg string) errorImpl {
	v := errorImpl{
		code: code,
	}

	if msg == "" {
		v.msg = strings.ReplaceAll(code, "_", " ")
	} else {
		v.msg = msg
	}

	return v
}

// notfoudError
type notfoudError struct {
	errorImpl
}

func (e notfoudError) NotFound() {}

// NewNotFound
func NewNotFound(code string, msg string) notfoudError {
	return notfoudError{New(code, msg)}
}

// noPermissionError
type noPermissionError struct {
	errorImpl
}

func (e noPermissionError) NoPermission() {}

// NewNoPermission
func NewNoPermission(msg string) noPermissionError {
	return noPermissionError{New(errorCodeNoPermission, msg)}
}
