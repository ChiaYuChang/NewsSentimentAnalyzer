package errorcode

type ErrorCode int32

// OK
const (
	Success          ErrorCode = 0
	SuccessNoContent ErrorCode = 1
)
const ECCodeHasBeenUsed ErrorCode = 2
const ECUnknownError ErrorCode = 3

// Client side general error
const (
	ECBadRequest ErrorCode = 400 + iota
	ECUnauthorized
	ECForbidden
	ECNotFound
	ECRequestTimeout
	ECTooManyRequests
	ECConflict
	ECUnsupportedMediaType
	ECUnprocessableContent
	ECInvalidParams
)

// Server side error
const (
	ECServerError ErrorCode = 500 + iota
	ECServiceUnavailable
)

// PGX error
const (
	ECPgxError ErrorCode = 600 + iota
)
