package httpresponse

import (
	"encoding/json"
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/errors"
)

var (
	defaultErrorResp []byte
)

func init() {
	defaultError := map[string]interface{}{
		"errors": []*errors.Error{
			errors.GetDefaultError(http.StatusInternalServerError),
		},
	}
	defaultErrorResp, _ = json.Marshal(defaultError)
}

// WithData wraps and writes the data to the http.ResponseWriter.
func WithData(w http.ResponseWriter, data interface{}) {
	status := http.StatusOK
	resp := map[string]interface{}{
		"data": data,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		jsonResp = defaultErrorResp
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(jsonResp))
}

// WithObject writes the object to the http.ResponseWriter.
func WithObject(w http.ResponseWriter, object interface{}) {
	status := http.StatusOK

	jsonResp, err := json.Marshal(object)
	if err != nil {
		jsonResp = defaultErrorResp
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(jsonResp))
}

// WithError writes the error to the http.ResponseWriter.
func WithError(w http.ResponseWriter, httpStatusCode int, errs ...error) {
	status := httpStatusCode
	resp := map[string]interface{}{
		"errors": errs,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		jsonResp = defaultErrorResp
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(jsonResp))
}

// InternalServerError writes internal server error with error message.
func InternalServerError(w http.ResponseWriter) {
	WithError(
		w,
		http.StatusInternalServerError,
		errors.GetDefaultError(http.StatusInternalServerError),
	)
}

// BadRequest writes bad requestr with error message.
func BadRequest(w http.ResponseWriter) {
	WithError(
		w,
		http.StatusBadRequest,
		errors.GetDefaultError(http.StatusBadRequest),
	)
}
