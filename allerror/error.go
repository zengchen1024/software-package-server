package allerror

import "strings"

const (
	ErrorCodeParamNotSpec                   = "param_not_spec_file"
	ErrorCodeParamNotSRPM                   = "param_not_srpm_file"
	ErrorCodeParamUserNotFound              = "param_user_not_found"
	ErrorCodeParamTooManyCommitters         = "param_too_many_committers"
	ErrorCodeParamDuplicateCommitters       = "param_duplicate_committers"
	ErrorCodeParamMissingChekItemComment    = "param_missing_check_item_comment"
	ErrorCodeParamImporterMissingPlatformId = "param_importer_missing_platform_id"

	errorCodeNoPermission = "no_permission"

	ErrorCodeAccessTokenMissing = "access_token_missing"
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	ErrorCodeSensitiveContent = "sensitive_content"

	ErrorCodePkgExists   = "software_pkg_exists"
	ErrorCodePkgNotFound = "software_pkg_not_found"

	ErrorCodeCIIsRunning      = "software_pkg_ci_is_running"
	ErrorCodeCIIsNotReady     = "software_pkg_ci_is_not_ready"
	ErrorCodeRetestRepeatedly = "software_pkg_retest_repeatedly"

	ErrorCodePkgCodeNotReady   = "software_pkg_code_not_ready"
	ErrorCodePkgIncorrectPhase = "software_pkg_incorrect_phase"
	ErrorCodePkgNothingChanged = "software_pkg_nothing_changed"

	ErrorCodePkgCommentNotFound = "software_pkg_comment_not_found"
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
