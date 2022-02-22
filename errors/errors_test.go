package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New(500, "500", "Internal Server Error", "Internal Server Error")
	assert.NotNil(t, err)

	opErr := NewOpError("New", fmt.Errorf("Failed to create new error"))
	assert.NotNil(t, opErr)
}

func TestError_WithSource(t *testing.T) {
	err := New(500, "500", "Internal Server Error", "Internal Server Error")
	assert.NotNil(t, err)
	assert.Nil(t, err.Source)

	err = err.WithSource("/panic", "")
	assert.NotNil(t, err.Source)
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Status int
		Code   string
		Title  string
		Detail string
		Source *ErrorSource
		Op     string
		Err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Error for end user, but empty",
			fields: fields{},
			want:   `{"status":"","code":"","title":"","detail":""}`,
		},
		{
			name: "Error for end user",
			fields: fields{
				Status: 500,
				Code:   "500",
				Title:  "Internal Server Error",
				Detail: "Internal Server Error",
				Source: &ErrorSource{
					Header: "Authorization",
				},
			},
			want: `{"status":500,"code":"500","title":"Internal Server Error","detail":"Internal Server Error","source":{"parameter":"","header":"Authorization"}}`,
		},
		{
			name: "Error for operator",
			fields: fields{
				Op:  "Op",
				Err: fmt.Errorf("Err"),
			},
			want: `Op: Err`,
		},
		{
			name: "Nested error",
			fields: fields{
				Op: "OpParent",
				Err: &Error{
					Op:  "OpChild",
					Err: fmt.Errorf("ErrChild"),
				},
			},
			want: `OpParent: OpChild: ErrChild`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Status: tt.fields.Status,
				Code:   tt.fields.Code,
				Title:  tt.fields.Title,
				Detail: tt.fields.Detail,
				Source: tt.fields.Source,
				Op:     tt.fields.Op,
				Err:    tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ErrorDetail(t *testing.T) {
	type fields struct {
		Status int
		Code   string
		Title  string
		Detail string
		Source *ErrorSource
		Op     string
		Err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "no error",
			fields: fields{},
			want:   "An internal error has occurred.",
		},
		{
			name: "with detail",
			fields: fields{
				Detail: "Error detail.",
			},
			want: "Error detail.",
		},
		{
			name: "with error",
			fields: fields{
				Err: fmt.Errorf("Error"),
			},
			want: "Error",
		},
		{
			name: "error with nested message",
			fields: fields{
				Err: &Error{
					Detail: "Child's error detail.",
				},
			},
			want: "Child's error detail.",
		},
		{
			name: "error with nested error",
			fields: fields{
				Err: &Error{
					Err: fmt.Errorf("Child's error"),
				},
			},
			want: "Child's error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Status: tt.fields.Status,
				Code:   tt.fields.Code,
				Title:  tt.fields.Title,
				Detail: tt.fields.Detail,
				Source: tt.fields.Source,
				Op:     tt.fields.Op,
				Err:    tt.fields.Err,
			}
			if got := e.ErrorDetail(); got != tt.want {
				t.Errorf("Error.ErrorDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorDetail_nil(t *testing.T) {
	got := errorDetail(nil)
	assert.Equal(t, got, "")
}
