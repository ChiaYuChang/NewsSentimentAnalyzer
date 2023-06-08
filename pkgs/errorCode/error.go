package errorcode

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Error struct {
	ErrorCode      ErrorCode `json:"error_code"`
	HttpStatusCode int       `json:"status_code"`
	Message        string    `json:"message"`
	MessageF       string    `json:"-"`
	Details        []string  `json:"details,omitempty"`
}

func NewError(ec ErrorCode, status int, msg string) *Error {
	return &Error{
		ErrorCode:      ec,
		HttpStatusCode: status,
		Message:        msg,
		Details:        nil,
	}
}

func NewErrorFromErr(ec ErrorCode, status int, err error) *Error {
	return NewError(ec, status, err.Error())
}

func (e *Error) Clone() *Error {
	newErr := NewError(e.ErrorCode, e.HttpStatusCode, e.Message)
	if len(e.Details) > 0 {
		newErr.Details = make([]string, len(e.Details))
		copy(newErr.Details, e.Details)
	}
	return newErr
}

func (e Error) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("code: %d, status: %d, msg: %s", e.ErrorCode, e.HttpStatusCode, e.Message))

	if len(e.Details) > 0 {
		sb.WriteString(", details:\n")
		for i, d := range e.Details {
			sb.WriteString(fmt.Sprintf("\t- %2d: %s\n", i, d))
		}
	}
	return sb.String()
}

func (e *Error) WithMessage(msg string) *Error {
	e.Message = msg
	return e
}

func (e *Error) WithMessagef(msgf string) *Error {
	e.MessageF = msgf
	return e
}

func (e *Error) WithDetails(details ...string) *Error {
	e.Details = append(e.Details, details...)
	return e
}

func (e Error) Msg() string {
	return e.Message
}

func (e Error) Msgf(a ...any) string {
	return fmt.Sprintf(e.MessageF, a...)
}

func (e Error) IsEqual(err error) bool {
	if ecErr, ok := err.(*Error); ok {
		return e.ErrorCode == ecErr.ErrorCode
	}
	return false
}

func (e Error) ToJson() ([]byte, error) {
	return json.MarshalIndent(e, "", "    ")
}

func (e Error) MustToJson() []byte {
	b, _ := json.MarshalIndent(e, "", "    ")
	return b
}
