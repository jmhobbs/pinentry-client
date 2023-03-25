package pinentry

import (
	"fmt"

	"github.com/jmhobbs/pinentry-client/assuan"
)

const (
	ErrorCodeCancelled    string = "83886179"
	ErrorCodeNotConfirmed        = "83886194"
)

type PinentryError struct {
	Code        string
	Description string
}

func (e PinentryError) Error() string {
	return fmt.Sprintf("%s (%s)", e.Description, e.Code)
}

type PinentryCancelledError struct {
	PinentryError
}

type PinentryNotConfirmedError struct {
	PinentryError
}

func NewPinentryError(response assuan.Response) error {
	switch response.Code {
	case ErrorCodeCancelled:
		return PinentryCancelledError{PinentryError{Description: response.Description}}
	case ErrorCodeNotConfirmed:
		return PinentryNotConfirmedError{PinentryError{Description: response.Description}}
	default:
		return PinentryError{Code: response.Code, Description: response.Description}
	}
}
