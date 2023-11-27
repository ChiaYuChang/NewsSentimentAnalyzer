package errorcode

import "net/http"

type ErrorCode int32

// OK
const (
	Success          ErrorCode = 0
	SuccessNoContent ErrorCode = 1
)

const (
	ECCodeHasBeenUsed ErrorCode = 2
	ECUnknownError    ErrorCode = 3
	ECgRPCServerError ErrorCode = 4
	ECgRPCClientError ErrorCode = 5
)

// Client side general error
// HTTP status codes
const (
	ECBadRequest           ErrorCode = http.StatusBadRequest
	ECUnauthorized         ErrorCode = http.StatusUnauthorized
	ECForbidden            ErrorCode = http.StatusForbidden
	ECNotFound             ErrorCode = http.StatusNotFound
	ECRequestTimeout       ErrorCode = http.StatusRequestTimeout
	ECTooManyRequests      ErrorCode = http.StatusTooManyRequests
	ECConflict             ErrorCode = http.StatusConflict
	ECUnsupportedMediaType ErrorCode = http.StatusUnsupportedMediaType
	ECUnprocessableContent ErrorCode = http.StatusUnprocessableEntity
	ECGone                 ErrorCode = http.StatusGone
)

// Self defined client error
const (
	ECInvalidParams ErrorCode = 460
)

// Server side error
const (
	ECServerError        ErrorCode = http.StatusInternalServerError
	ECServiceUnavailable ErrorCode = http.StatusServiceUnavailable
)

// PGX error
const (
	ECPgxError ErrorCode = 600 + iota
)
