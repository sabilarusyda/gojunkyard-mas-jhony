package wsresponse

import (
	"reflect"

	"devcode.xeemore.com/systech/gojunkyard/websocket/connection"
)

type response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

// JSONWithData wraps and writes the data to the websocket connection.
func JSONWithData(conn *connection.Connection, data interface{}) {
	if !reflect.ValueOf(data).IsNil() {
		conn.WriteJSON(response{
			Status: "OK",
			Data:   data,
		})
	}
}

// JSONWithError writes the error to the websocket connection.
func JSONWithError(conn *connection.Connection, err error) {
	conn.WriteJSON(response{
		Status: "FAILED",
		Error:  err.Error(),
	})
}
