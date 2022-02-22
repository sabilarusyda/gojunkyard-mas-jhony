package valkyrie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestHTTPRouter_AuthProject(t *testing.T) {
	type fields struct {
		core *core
	}
	type args struct {
		r *http.Request
	}
	type want struct {
		pid int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "without project",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
		},
		{
			name: "with invalid project",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/test?app_id=abd", nil),
			},
		},
		{
			name: "with valid project (querystring)",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/test?app_id=molatv", nil),
			},
			want: want{
				pid: 2,
			},
		},
		{
			name: "with valid project (x-url)",
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/test", nil)
					r.Header.Set("x-url", "analytic.supersoccer.tv")
					return r
				}(),
			},
			want: want{
				pid: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				c  = &HTTPRouter{getTestCore()}
				w  = httptest.NewRecorder()
				r  = tt.args.r
				ps = httprouter.Params{}
			)

			f := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				pid, _ := GetProjectID(r)
				if pid != tt.want.pid {
					t.Errorf("HTTPRouter().pid = %v, want %v", pid, tt.want.pid)
				}
			}

			got := c.AuthProject(f)
			got(w, r, ps)
		})
	}
}
