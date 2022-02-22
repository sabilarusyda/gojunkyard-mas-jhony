package errors

import (
	"fmt"
	"strings"
)

// Error defines a standard application error.
type Error struct {
	// For end user.
	// an HTTP status code value, without the textual description
	Status int `json:"status"`
	// an application-specific error code
	Code string `json:"code"`
	// a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem
	Title string `json:"title"`
	// a human-readable explanation specific to this occurrence of the problem
	Detail string `json:"detail"`
	// an object containing references to the source of the error
	Source *ErrorSource `json:"source,omitempty"`

	// Logical operation and nested error. For operator / developer.
	Op  string `json:"-"`
	Err error  `json:"-"`
}

// ErrorSource references to the source of the error.
type ErrorSource struct {
	Parameter string `json:"parameter,omitempty"` // a string indicating which URI query parameter caused the error
	Header    string `json:"header,omitempty"`    // a string indicating which request header caused the error
}

// New returns an initialized error for end-user.
func New(status int, code, title, detail string) *Error {
	return &Error{
		Status: status,
		Code:   code,
		Title:  title,
		Detail: detail,
	}
}

// NewOpError returns an initialized error for operator / developer.
func NewOpError(op string, err error) *Error {
	return &Error{
		Op:  op,
		Err: err,
	}
}

// WithSource returns an error with error source.
func (e *Error) WithSource(parameter, header string) *Error {
	clone := e.clone()
	clone.Source = &ErrorSource{
		Parameter: parameter,
		Header:    header,
	}
	return clone
}

func (e *Error) clone() *Error {
	clone := *e
	return &clone
}

// Error returns the string representation of the error.
// If Err exists, we only write Err and Op (if any).
func (e *Error) Error() string {
	var buf strings.Builder

	// If wrapping an error, print its Error() message.
	// Otherwise print anything else.
	if e.Err != nil {
		// Print the current operation in our stack, if any.
		if e.Op != "" {
			fmt.Fprintf(&buf, "%s: ", e.Op)
		}
		buf.WriteString(e.Err.Error())
	} else {
		fmt.Fprintf(
			&buf,
			`{"status":"%d","code":"%s","title":"%s","detail":"%s"`,
			e.Status, e.Code, e.Title, e.Detail,
		)
		if e.Source != nil {
			fmt.Fprintf(
				&buf,
				`,"source":{"parameter":"%s","header":"%s"}`,
				e.Source.Parameter, e.Source.Header,
			)
		}
		buf.WriteRune('}')
	}

	return buf.String()
}

// ErrorDetail returns the human-readable message detail of the error, if available.
// Otherwise returns a generic error message.
func (e *Error) ErrorDetail() string {
	return errorDetail(e)
}

func errorDetail(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Detail != "" {
		return e.Detail
	} else if ok && e.Err != nil {
		return errorDetail(e.Err)
	} else if !ok {
		return err.Error()
	}
	return "An internal error has occurred."
}
