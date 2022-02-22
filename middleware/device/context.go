package device

import "net/http"

type ctxDevice uint8

const ctxDeviceID ctxDevice = iota + 1

// IsDevice ...
func IsDevice(r *http.Request) bool {
	_, ok := r.Context().Value(ctxDeviceID).(string)
	return ok
}

// GetDeviceID ...
func GetDeviceID(r *http.Request) string {
	str, _ := r.Context().Value(ctxDeviceID).(string)
	return str
}
