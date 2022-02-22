package device

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type clientMock struct {
	mock.Mock
}

func (cm *clientMock) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	args := cm.Called(in)
	return args.Get(0).(*ValidateResponse), args.Error(1)
}

func TestDevice_Close(t *testing.T) {
	type fields struct {
		client ValidatorServiceClient
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Empty conn",
			fields: fields{
				client: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				client: tt.fields.client,
			}
			if err := d.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Device.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDevice_HandleOptional(t *testing.T) {
	type fields struct {
		client ValidatorServiceClient
	}
	type args struct {
		r *http.Request
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "No TOTP",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			want: want{
				code: http.StatusOK,
				body: "OK",
			},
		},
		{
			name: "Client Error",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req.Header.Add("x-totp", "TOTP_VALUE")
					req.Header.Add("user-agent", "httpmock")
					req.Header.Add("x-forwarded-for", "172.16.0.2")
					return req
				}(),
			},
			fields: fields{
				client: func() *clientMock {
					cm := new(clientMock)
					cm.On("Validate", &ValidateRequest{
						Totp: "TOTP_VALUE",
						Ua:   "httpmock",
						Ip:   "172.16.0.2",
					}).Return((*ValidateResponse)(nil), errors.New("Error connection"))
					return cm
				}(),
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "{\"errors\":[{\"status\":\"500\",\"code\":\"00900001\",\"title\":\"INTERNAL_SERVER_ERROR\",\"detail\":\"INTERNAL_SERVER_ERROR\"}]}",
			},
		},
		{
			name: "Invalid TOTP",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req.Header.Add("x-totp", "TOTP_VALUE")
					req.Header.Add("user-agent", "httpmock")
					req.Header.Add("x-forwarded-for", "172.16.0.2")
					return req
				}(),
			},
			fields: fields{
				client: func() *clientMock {
					cm := new(clientMock)
					cm.On("Validate", &ValidateRequest{
						Totp: "TOTP_VALUE",
						Ua:   "httpmock",
						Ip:   "172.16.0.2",
					}).Return(&ValidateResponse{
						Valid: false,
					}, (error)(nil))
					return cm
				}(),
			},
			want: want{
				code: http.StatusForbidden,
				body: "{\"errors\":[{\"status\":\"403\",\"code\":\"00900002\",\"title\":\"FORBIDDEN\",\"detail\":\"INVALID_TOKEN\"}]}",
			},
		},
		{
			name: "Valid TOTP",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req.Header.Add("x-totp", "TOTP_VALUE")
					req.Header.Add("user-agent", "httpmock")
					req.Header.Add("x-forwarded-for", "172.16.0.2")
					return req
				}(),
			},
			fields: fields{
				client: func() *clientMock {
					cm := new(clientMock)
					cm.On("Validate", &ValidateRequest{
						Totp: "TOTP_VALUE",
						Ua:   "httpmock",
						Ip:   "172.16.0.2",
					}).Return(&ValidateResponse{
						Valid:    true,
						DeviceId: "1ab822ab8d3ee1ab822ab8d3ee1ab822",
					}, (error)(nil))
					return cm
				}(),
			},
			want: want{
				code: http.StatusOK,
				body: "OK",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{client: tt.fields.client}

			// === HandleOptional === //
			w0 := httptest.NewRecorder()
			d.HandleOptional(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})).ServeHTTP(w0, tt.args.r)

			assert.Equal(t, tt.want.body, w0.Body.String())
			assert.Equal(t, tt.want.code, w0.Code)

			// === HandleMust === //
			w1 := httptest.NewRecorder()
			d.HandleFuncOptional(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			}).ServeHTTP(w1, tt.args.r)

			assert.Equal(t, tt.want.body, w1.Body.String())
			assert.Equal(t, tt.want.code, w1.Code)
		})
	}
}

func TestDevice_HandleMust(t *testing.T) {
	type fields struct {
		client ValidatorServiceClient
	}
	type args struct {
		r *http.Request
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "No TOTP",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			want: want{
				code: http.StatusForbidden,
				body: "{\"errors\":[{\"status\":\"403\",\"code\":\"00900002\",\"title\":\"FORBIDDEN\",\"detail\":\"INVALID_TOKEN\"}]}",
			},
		},
		{
			name: "Client Error",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req.Header.Add("x-totp", "TOTP_VALUE")
					req.Header.Add("user-agent", "httpmock")
					req.Header.Add("x-forwarded-for", "172.16.0.2")
					return req
				}(),
			},
			fields: fields{
				client: func() *clientMock {
					cm := new(clientMock)
					cm.On("Validate", &ValidateRequest{
						Totp: "TOTP_VALUE",
						Ua:   "httpmock",
						Ip:   "172.16.0.2",
					}).Return((*ValidateResponse)(nil), errors.New("Error connection"))
					return cm
				}(),
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "{\"errors\":[{\"status\":\"500\",\"code\":\"00900001\",\"title\":\"INTERNAL_SERVER_ERROR\",\"detail\":\"INTERNAL_SERVER_ERROR\"}]}",
			},
		},
		{
			name: "Invalid TOTP",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req.Header.Add("x-totp", "TOTP_VALUE")
					req.Header.Add("user-agent", "httpmock")
					req.Header.Add("x-forwarded-for", "172.16.0.2")
					return req
				}(),
			},
			fields: fields{
				client: func() *clientMock {
					cm := new(clientMock)
					cm.On("Validate", &ValidateRequest{
						Totp: "TOTP_VALUE",
						Ua:   "httpmock",
						Ip:   "172.16.0.2",
					}).Return(&ValidateResponse{
						Valid: false,
					}, (error)(nil))
					return cm
				}(),
			},
			want: want{
				code: http.StatusForbidden,
				body: "{\"errors\":[{\"status\":\"403\",\"code\":\"00900002\",\"title\":\"FORBIDDEN\",\"detail\":\"INVALID_TOKEN\"}]}",
			},
		},
		{
			name: "Valid TOTP",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/", nil)
					req.Header.Add("x-totp", "TOTP_VALUE")
					req.Header.Add("user-agent", "httpmock")
					req.Header.Add("x-forwarded-for", "172.16.0.2")
					return req
				}(),
			},
			fields: fields{
				client: func() *clientMock {
					cm := new(clientMock)
					cm.On("Validate", &ValidateRequest{
						Totp: "TOTP_VALUE",
						Ua:   "httpmock",
						Ip:   "172.16.0.2",
					}).Return(&ValidateResponse{
						Valid:    true,
						DeviceId: "1ab822ab8d3ee1ab822ab8d3ee1ab822",
					}, (error)(nil))
					return cm
				}(),
			},
			want: want{
				code: http.StatusOK,
				body: "OK",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{client: tt.fields.client}

			// === HandleOptional === //
			w0 := httptest.NewRecorder()
			d.HandleMust(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})).ServeHTTP(w0, tt.args.r)

			assert.Equal(t, tt.want.body, w0.Body.String())
			assert.Equal(t, tt.want.code, w0.Code)

			// === HandleMust === //
			w1 := httptest.NewRecorder()
			d.HandleFuncMust(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			}).ServeHTTP(w1, tt.args.r)

			assert.Equal(t, tt.want.body, w1.Body.String())
			assert.Equal(t, tt.want.code, w1.Code)
		})
	}
}
