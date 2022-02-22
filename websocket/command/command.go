package command

import (
	"encoding/json"
	"errors"

	"devcode.xeemore.com/systech/gojunkyard/websocket/connection"
)

var commandFunc = make(map[string]Command, 20)

// Request ...
type Request struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Response ...
type Response struct {
	ID         interface{} `json:"id"`
	Type       string      `json:"type"`
	Attributes interface{} `json:"attributes,omitempty"`
	Error      error       `json:"error,omitempty"`
}

// Command ...
type Command interface {
	Name() string
	Exec(conn *connection.Connection, message []byte) (v *Response, err error)
}

// Register ...
func Register(cmds []Command) {
	for _, v := range cmds {
		commandFunc[v.Name()] = v
	}
}

// Run ...
func Run(conn *connection.Connection, message []byte) (v *Response, err error) {
	var data json.RawMessage
	schema := Request{
		Data: &data,
	}

	if err := json.Unmarshal(message, &schema); err != nil {
		return nil, err
	}

	c, ok := commandFunc[schema.Type]
	if !ok {
		return nil, errors.New("Invalid Command")
	}
	return c.Exec(conn, data)
}
