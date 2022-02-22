package errors

import (
	"net/http"
	"strconv"
)

var errorCodes = []int{
	http.StatusBadRequest,                    // 400 - RFC 7231, 6.5.1
	http.StatusUnauthorized,                  // 401 - RFC 7235, 3.1
	http.StatusPaymentRequired,               // 402 - RFC 7231, 6.5.2
	http.StatusForbidden,                     // 403 - RFC 7231, 6.5.3
	http.StatusNotFound,                      // 404 - RFC 7231, 6.5.4
	http.StatusMethodNotAllowed,              // 405 - RFC 7231, 6.5.5
	http.StatusNotAcceptable,                 // 406 - RFC 7231, 6.5.6
	http.StatusProxyAuthRequired,             // 407 - RFC 7235, 3.2
	http.StatusRequestTimeout,                // 408 - RFC 7231, 6.5.7
	http.StatusConflict,                      // 409 - RFC 7231, 6.5.8
	http.StatusGone,                          // 410 - RFC 7231, 6.5.9
	http.StatusLengthRequired,                // 411 - RFC 7231, 6.5.10
	http.StatusPreconditionFailed,            // 412 - RFC 7232, 4.2
	http.StatusRequestEntityTooLarge,         // 413 - RFC 7231, 6.5.11
	http.StatusRequestURITooLong,             // 414 - RFC 7231, 6.5.12
	http.StatusUnsupportedMediaType,          // 415 - RFC 7231, 6.5.13
	http.StatusRequestedRangeNotSatisfiable,  // 416 - RFC 7233, 4.4
	http.StatusExpectationFailed,             // 417 - RFC 7231, 6.5.14
	http.StatusTeapot,                        // 418 - RFC 7168, 2.3.3
	http.StatusUnprocessableEntity,           // 422 - RFC 4918, 11.2
	http.StatusLocked,                        // 423 - RFC 4918, 11.3
	http.StatusFailedDependency,              // 424 - RFC 4918, 11.4
	http.StatusUpgradeRequired,               // 426 - RFC 7231, 6.5.15
	http.StatusPreconditionRequired,          // 428 - RFC 6585, 3
	http.StatusTooManyRequests,               // 429 - RFC 6585, 4
	http.StatusRequestHeaderFieldsTooLarge,   // 431 - RFC 6585, 5
	http.StatusUnavailableForLegalReasons,    // 451 - RFC 7725, 3
	http.StatusInternalServerError,           // 500 - RFC 7231, 6.6.1
	http.StatusNotImplemented,                // 501 - RFC 7231, 6.6.2
	http.StatusBadGateway,                    // 502 - RFC 7231, 6.6.3
	http.StatusServiceUnavailable,            // 503 - RFC 7231, 6.6.4
	http.StatusGatewayTimeout,                // 504 - RFC 7231, 6.6.5
	http.StatusHTTPVersionNotSupported,       // 505 - RFC 7231, 6.6.6
	http.StatusVariantAlsoNegotiates,         // 506 - RFC 2295, 8.1
	http.StatusInsufficientStorage,           // 507 - RFC 4918, 11.5
	http.StatusLoopDetected,                  // 508 - RFC 5842, 7.2
	http.StatusNotExtended,                   // 510 - RFC 2774, 7
	http.StatusNetworkAuthenticationRequired, // 511 - RFC 6585, 6
}

// defaultErrors maps http status code to it's error
var defaultErrors map[int]*Error

func init() {
	defaultErrors = make(map[int]*Error)
	for _, errorCode := range errorCodes {
		defaultErrors[errorCode] = &Error{
			Status: errorCode,
			Code:   strconv.FormatInt(int64(errorCode), 10),
			Title:  http.StatusText(errorCode),
			Detail: http.StatusText(errorCode),
		}
	}
}

// GetDefaultError returns default error that wrap http status code and text.
func GetDefaultError(code int) *Error {
	return defaultErrors[code]
}
