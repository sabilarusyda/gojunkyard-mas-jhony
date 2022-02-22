package device

import (
	context "context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsDevice(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Not Device",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			want: false,
		},
		{
			name: "True Device",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req = req.WithContext(context.WithValue(req.Context(), ctxDeviceID, "abc"))
					return req
				}(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDevice(tt.args.r); got != tt.want {
				t.Errorf("IsDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDeviceID(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Not Device",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			want: "",
		},
		{
			name: "True Device",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req = req.WithContext(context.WithValue(req.Context(), ctxDeviceID, "abc"))
					return req
				}(),
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDeviceID(tt.args.r); got != tt.want {
				t.Errorf("GetDeviceID() = %v, want %v", got, tt.want)
			}
		})
	}
}
