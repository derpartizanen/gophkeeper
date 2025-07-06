package errors

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const UnknownErrorCode = 255

// RequestError is custom error wrapper used inside keeperctl to distinguish
// internal errors from returned from keeperd service.
type RequestError struct {
	code        uint32
	description string
	details     []string
}

// NewRequestError create RequestError from provided possibly gRPC error.
func NewRequestError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return RequestError{
			code:        UnknownErrorCode,
			description: err.Error(),
		}
	}

	details := make([]string, 0)

	for _, detail := range st.Details() {
		if t, ok := detail.(*errdetails.BadRequest); ok {
			for _, violation := range t.GetFieldViolations() {
				details = append(
					details,
					fmt.Sprintf("%q: %s", violation.GetField(), violation.GetDescription()),
				)
			}
		}
	}

	return RequestError{
		code:        uint32(st.Code()),
		description: st.Message(),
		details:     details,
	}
}

// Error returns error text.
func (e RequestError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s (%d)", e.description, e.code))

	for _, detail := range e.details {
		sb.WriteString(fmt.Sprintf("\n\t%s", detail))
	}

	return sb.String()
}

// Unwrap get error and unwraps it to RequestError.
func Unwrap(err error) error {
	var rErr RequestError

	if errors.As(err, &rErr) {
		return rErr
	}

	return err
}
