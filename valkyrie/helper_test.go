package valkyrie

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProjectID(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type want struct {
		pid   int64
		valid bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "forget to add middleware",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/test", nil),
			},
		},
		{
			name: "success",
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/test?app_id=supersoccertv", nil)
					r = r.WithContext(context.WithValue(r.Context(), PID, int64(1)))
					return r
				}(),
			},
			want: want{
				pid:   1,
				valid: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPid, gotValid := GetProjectID(tt.args.r)
			if gotPid != tt.want.pid {
				t.Errorf("GetProjectID() gotPid = %v, want %v", gotPid, tt.want.pid)
			}
			if gotValid != tt.want.valid {
				t.Errorf("GetProjectID() gotValid = %v, want %v", gotValid, tt.want.valid)
			}
		})
	}
}
