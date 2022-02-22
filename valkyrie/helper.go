package valkyrie

import "net/http"

// GetProjectID is the helper function which is used to get project id from request.
// To make this function work. The middleware is required
func GetProjectID(r *http.Request) (pid int64, valid bool) {
	pid, valid = r.Context().Value(PID).(int64)
	return
}
