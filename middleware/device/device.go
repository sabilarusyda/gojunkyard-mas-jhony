package device

import (
	context "context"
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/errors"
	"devcode.xeemore.com/systech/gojunkyard/http/httpresponse"

	"google.golang.org/grpc"
)

// Device ...
type Device struct {
	client ValidatorServiceClient
}

var (
	// ErrConnection happens when failed connecting to totp validator server
	ErrConnection = errors.New(500, "00900001", "INTERNAL_SERVER_ERROR", "INTERNAL_SERVER_ERROR")
	// ErrInvalidToken is returned when totp is invalid
	ErrInvalidToken = errors.New(403, "00900002", "FORBIDDEN", "INVALID_TOKEN")
)

// New ...
// Example addr: validator.accounts.svc.local:9000
func New(addr string) (*Device, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Device{
		client: NewValidatorServiceClient(conn),
	}, nil
}

// HandleFuncMust ...
func (d *Device) HandleFuncMust(h http.HandlerFunc) http.HandlerFunc {
	return d.HandleMust(h).(http.HandlerFunc)
}

// HandleMust ...
func (d *Device) HandleMust(h http.Handler) http.Handler {
	return d.HandleOptional(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsDevice(r) {
			httpresponse.WithError(w, http.StatusForbidden, ErrInvalidToken)
			return
		}

		h.ServeHTTP(w, r)
	}))
}

// HandleFuncOptional ...
func (d *Device) HandleFuncOptional(h http.HandlerFunc) http.HandlerFunc {
	return d.HandleOptional(h).(http.HandlerFunc)
}

// HandleOptional ...
func (d *Device) HandleOptional(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xTotp := r.Header.Get("x-totp")
		if len(xTotp) == 0 {
			h.ServeHTTP(w, r)
			return
		}

		ip := r.RemoteAddr
		if xForwardedFor := r.Header.Get("x-forwarded-for"); len(xForwardedFor) > 0 {
			ip = xForwardedFor
		}

		resp, err := d.client.Validate(r.Context(), &ValidateRequest{
			Totp: xTotp,
			Ip:   ip,
			Ua:   r.Header.Get("user-agent"),
		})
		if err != nil {
			httpresponse.WithError(w, http.StatusInternalServerError, ErrConnection)
			return
		}

		if !resp.Valid {
			httpresponse.WithError(w, http.StatusForbidden, ErrInvalidToken)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxDeviceID, resp.DeviceId)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Close ...
func (d *Device) Close() error {
	client, _ := d.client.(*validatorServiceClient)
	if client == nil {
		return nil
	}
	return client.cc.Close()
}
