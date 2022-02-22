package valkyrie

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func getTestOption() *Option {
	return &Option{
		URL:           "https://staging.supersoccer.tv/api/v2/config/apps",
		CacheDuration: time.Minute,
	}
}

func getTestCore() *core {
	c := &core{
		client: new(http200Client),
		option: getTestOption(),
	}
	c.loadHostmap()
	return c
}

func Test_examineOption(t *testing.T) {
	type mock struct {
		logFatalf func(format string, v ...interface{})
	}
	type args struct {
		option *Option
	}
	type want struct {
		option *Option
	}
	tests := []struct {
		name string
		args args
		mock mock
		want want
	}{
		{
			name: "nil option",
			mock: mock{
				logFatalf: func(format string, v ...interface{}) {
					const want = "Failed to initialize valkyrie. option cannot be nil"
					if format != want {
						log.Fatalf("got: %s, want: %s", format, want)
					}
					if v != nil {
						log.Fatalf("got: %v, want: %v", v, nil)
					}
				},
			},
		},
		{
			name: "empty url",
			args: args{
				option: &Option{},
			},
			mock: mock{
				logFatalf: func(format string, v ...interface{}) {
					const want = "Failed to initialize valkyrie. option.url cannot be empty"
					if format != want {
						log.Fatalf("got: %s, want: %s", format, want)
					}
					if v != nil {
						log.Fatalf("got: %v, want: %v", v, nil)
					}
				},
			},
			want: want{
				option: &Option{},
			},
		},
		{
			name: "invalid url format",
			args: args{
				option: &Option{
					URL: "a.c.v.s.sdl.asdfk123",
				},
			},
			mock: mock{
				logFatalf: func(format string, v ...interface{}) {
					const want = "Failed to initialize valkyrie. option.url got invalid format. err: %s"
					const v0 = `parse "a.c.v.s.sdl.asdfk123": invalid URI for request`
					if format != want {
						log.Fatalf("got: %s, want: %s", format, want)
					}
					if len(v) != 1 {
						log.Fatalf("got.len: %d, want: 1", len(v))
					}
					err, ok := v[0].(error)
					if !ok {
						log.Fatalf("v0 must be an error")
					}
					if err.Error() != v0 {
						log.Fatalf("got: %v, want: %v", v, v0)
					}
				},
			},
			want: want{
				option: &Option{
					URL: "a.c.v.s.sdl.asdfk123",
				},
			},
		},
		{
			name: "cache duration 0",
			args: args{
				option: &Option{
					CacheDuration: 0,
					URL:           "https://staging.supersoccer.tv/api/v2/config/apps",
				},
			},
			want: want{
				option: &Option{
					CacheDuration: time.Minute,
					URL:           "https://staging.supersoccer.tv/api/v2/config/apps",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFatalf = tt.mock.logFatalf
			examineOption(tt.args.option)
			if !reflect.DeepEqual(tt.args.option, tt.want.option) {
				t.Errorf("got: %v, want: %v", tt.args.option, tt.want.option)
			}
		})
	}
}

func Test_core_loadHostmap(t *testing.T) {
	tests := []struct {
		name    string
		c       *core
		wantErr bool
	}{
		{
			name: "cannot connect to server",
			c: &core{
				option: getTestOption(),
				client: new(httpFailClient),
			},
			wantErr: true,
		},
		{
			name: "internal server error",
			c: &core{
				option: getTestOption(),
				client: new(http500Client),
			},
			wantErr: true,
		},
		{
			name: "status ok, but nil body",
			c: &core{
				option: getTestOption(),
				client: new(http200NilClient),
			},
			wantErr: true,
		},
		{
			name: "status ok, response valid certs",
			c: &core{
				option: getTestOption(),
				client: new(http200Client),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.loadHostmap(); (err != nil) != tt.wantErr {
				t.Errorf("core.loadHostmap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_core_authProject(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type want struct {
		pid int64
		ok  bool
	}
	tests := []struct {
		name string
		c    *core
		args args
		want want
	}{
		{
			name: "empty identifier",
			c:    getTestCore(),
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/test", nil),
			},
		},
		{
			name: "valid identifier. cached",
			c:    getTestCore(),
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/test?app_id=supersoccertv_ads", nil),
			},
			want: want{
				ok:  true,
				pid: 1,
			},
		},
		{
			name: "invalid identifier",
			c:    getTestCore(),
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/test?app_id=supermantaptv", nil),
			},
		},
		{
			name: "valid identifier. request new hostmap",
			c: &core{
				client: new(http200Client),
				option: getTestOption(),
			},
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/test?app_id=supersoccertv_ads", nil),
			},
			want: want{
				ok:  true,
				pid: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPid, gotOk := tt.c.authProject(tt.args.r)
			if gotPid != tt.want.pid {
				t.Errorf("core.authProjectID() gotPid = %v, want %v", gotPid, tt.want.pid)
			}
			if gotOk != tt.want.ok {
				t.Errorf("core.authProjectID() gotOk = %v, want %v", gotOk, tt.want.ok)
			}
		})
	}
}

func Test_core_authHostname(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type want struct {
		pid int64
		ok  bool
	}
	tests := []struct {
		name string
		c    *core
		args args
		want want
	}{
		{
			name: "valid hostname",
			c:    getTestCore(),
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/test", nil)
					r.Header.Add("x-url", "misty.mola.id/back-office")
					return r
				}(),
			},
			want: want{
				pid: 2,
				ok:  true,
			},
		},
		{
			name: "invalid hostname",
			c:    getTestCore(),
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/test", nil)
					r.Header.Add("x-url", "staging.supermantap.tv")
					return r
				}(),
			},
		},
		{
			name: "valid hostname, has not been cached",
			c: &core{
				client: new(http200Client),
				option: getTestOption(),
			},
			args: args{
				r: func() *http.Request {
					r := httptest.NewRequest(http.MethodGet, "/test", nil)
					r.Header.Add("x-url", "misty.mola.id/back-office")
					return r
				}(),
			},
			want: want{
				pid: 2,
				ok:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPid, gotOk := tt.c.authHostname(tt.args.r)
			if gotPid != tt.want.pid {
				t.Errorf("core.authHostname() gotPid = %v, want %v", gotPid, tt.want.pid)
			}
			if gotOk != tt.want.ok {
				t.Errorf("core.authHostname() gotOk = %v, want %v", gotOk, tt.want.ok)
			}
		})
	}
}

func Test_core_renewHostmap(t *testing.T) {
	c := &core{
		client: new(http200Client),
		option: &Option{
			URL:           "https://staging.supersoccer.tv/api/v2/config/apps",
			CacheDuration: 5 * time.Millisecond,
		},
	}
	if c.renewHostmap() != nil {
		t.Fatal("step 1. renew hostmap should be not error")
	}
	if c.renewHostmap() == nil {
		t.Fatal("step 2. renew hostmap should be error because of cache")
	}
}
