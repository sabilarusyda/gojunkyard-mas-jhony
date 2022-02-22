package httpresponse

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"devcode.xeemore.com/systech/gojunkyard/errors"

	"github.com/stretchr/testify/assert"
)

type emptyObject struct{}
type normalObject struct {
	Key1 string `json:"key_1"`
}

func compareResult(t *testing.T, want string, got string) {
	var wantObj, gotObj interface{}
	json.Unmarshal([]byte(want), &wantObj)
	json.Unmarshal([]byte(got), &gotObj)

	assert.Equal(t, wantObj, gotObj)
}

func TestWithData(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want string
	}{
		{
			name: "nil",
			data: nil,
			want: `{"data":null}`,
		},
		{
			name: "string",
			data: "supersoccer",
			want: `{"data":"supersoccer"}`,
		},
		{
			name: "empty object",
			data: emptyObject{},
			want: `{"data":{}}`,
		},
		{
			name: "normal object",
			data: normalObject{
				Key1: "value_1",
			},
			want: `{"data":{"key_1":"value_1"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			WithData(w, tt.data)

			compareResult(t, tt.want, w.Body.String())
		})
	}
}

func TestWithObject(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want string
	}{
		{
			name: "nil",
			data: nil,
			want: `null`,
		},
		{
			name: "string",
			data: "supersoccer",
			want: `"supersoccer"`,
		},
		{
			name: "empty object",
			data: emptyObject{},
			want: `{}`,
		},
		{
			name: "normal object",
			data: normalObject{
				Key1: "value_1",
			},
			want: `{"key_1":"value_1"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			WithObject(w, tt.data)

			compareResult(t, tt.want, w.Body.String())
		})
	}
}

func TestWithError(t *testing.T) {
	type args struct {
		httpStatusCode int
		errs           []error
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		want       string
	}{
		{
			name: "nil",
			args: args{
				httpStatusCode: 500,
				errs:           nil,
			},
			wantStatus: 500,
			want:       `{"errors":null}`,
		},
		{
			name: "one empty error",
			args: args{
				httpStatusCode: 500,
				errs: []error{
					&errors.Error{},
				},
			},
			wantStatus: 500,
			want:       `{"errors":[{"status":"","code":"","title":"","detail":""}]}`,
		},
		{
			name: "one error",
			args: args{
				httpStatusCode: 500,
				errs: []error{
					&errors.Error{
						Status: 500,
						Code:   "500",
						Title:  "Internal Server Error",
						Detail: "Internal Server Error",
					},
				},
			},
			wantStatus: 500,
			want:       `{"errors":[{"status":"500","code":"500","title":"Internal Server Error","detail":"Internal Server Error"}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			WithError(w, tt.args.httpStatusCode, tt.args.errs...)

			compareResult(t, fmt.Sprint(tt.wantStatus), fmt.Sprint(w.Code))
			compareResult(t, tt.want, w.Body.String())
		})
	}
}

func TestInternalServerError(t *testing.T) {
	want := `{"errors":[{"status":"500","code":"500","title":"Internal Server Error","detail":"Internal Server Error"}]}`

	w := httptest.NewRecorder()
	InternalServerError(w)

	compareResult(t, want, w.Body.String())
}

func TestBadRequest(t *testing.T) {
	want := `{"errors":[{"status":"400","code":"400","title":"Bad Request","detail":"Bad Request"}]}`

	w := httptest.NewRecorder()
	BadRequest(w)

	compareResult(t, want, w.Body.String())
}
