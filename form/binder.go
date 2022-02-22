package form

import (
	"context"
	"errors"
	"net/http"
	"strings"

	form "github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4/modifiers"
	validator "github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

type binderFlag uint8

const (
	// Bnone only does a form binding
	Bnone binderFlag = 0
	// Bfilter performs form binding and filter
	Bfilter binderFlag = 1 << iota
	// Bvalidate performs form binding and validation
	Bvalidate
	// Bstd performs form binding, filter, and validation
	Bstd = Bfilter | Bvalidate
)

var (
	_decoder   = form.NewDecoder()
	_filter    = modifiers.New()
	_validator = validator.New()
	// ErrNilRequest is returned when nil *request param is passed
	ErrNilRequest = errors.New("request cannot be nil")
	// ErrUnsupportedContentType is returned when content type parser is not exist
	ErrUnsupportedContentType = errors.New("content type is not supported")
	// jsoniterConfig is used for the configuration of json decoder
	jsonapi = jsoniter.Config{
		TagKey:                        "form",
		EscapeHTML:                    false,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
	}.Froze()
)

// Bind performs form binding, filter, and validation
/**
 * @param {interface{}} v - non nil pointer of any struct that will be binded
 * @param {*http.Request} r - non nil pointer of http.Request
 */
func Bind(v interface{}, r *http.Request) error {
	return BindFlag(Bstd, v, r)
}

// BindFlag performs form binding and optional (filter and validation)
/**
 * @param {flag} f - flag of action that will be perform (Bnone or Bfilter or Bvalidate or Bstd)
 * @param {interface{}} v - non nil pointer of any struct that will be binded
 * @param {*http.Request} r - non nil pointer of http.Request
 */
func BindFlag(f binderFlag, v interface{}, r *http.Request) error {
	// step 1. validate if r is nil, then return ErrNilRequest
	if r == nil {
		return ErrNilRequest
	}

	// step 2. decode request body based on request
	var mime = r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(mime, mimeJSON):
		// step 2a. decode json payload
		err := jsonapi.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			return err
		}
	case r.Method == http.MethodGet, strings.HasPrefix(mime, mimeMultipartForm), strings.HasPrefix(mime, mimeApplicationForm):
		// step 2b. decode form data payload
		err := r.ParseForm()
		if err != nil {
			return err
		}
		err = _decoder.Decode(v, r.Form)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedContentType
	}

	if f&Bfilter != 0 {
		_filter.Struct(context.Background(), v)
	}

	if f&Bvalidate != 0 {
		err := _validator.Struct(v)
		if err != nil {
			return err
		}
	}

	return nil
}
