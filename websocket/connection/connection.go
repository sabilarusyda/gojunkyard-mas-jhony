package connection

import (
	"context"
	"net/http"
	"sync"
	"time"

	"devcode.xeemore.com/systech/gojunkyard/util"

	"github.com/gorilla/websocket"
)

type Connection struct {
	id       string
	conn     *websocket.Conn
	ctx      context.Context
	channels channelSet
	close    func()
	mux      sync.RWMutex
	pid      int64
	appid    string
	room     string
	uid      string
}

type Option struct {
	HandshakeTimeout time.Duration
	ReadBufferSize   int
	WriteBufferSize  int
	Origins          []string
	OriginValid      bool
}

func New(w http.ResponseWriter, r *http.Request, option *Option) (*Connection, error) {
	var connectionUpgrader = websocket.Upgrader{
		HandshakeTimeout: option.HandshakeTimeout,
		ReadBufferSize:   option.ReadBufferSize,
		WriteBufferSize:  option.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return option.OriginValid
		},
	}
	conn, err := connectionUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	conn.SetCloseHandler(nil)
	conn.SetPingHandler(nil)
	conn.SetPongHandler(nil)
	return &Connection{id: util.GenerateUUID(), conn: conn}, nil
}

func (c *Connection) SetPongHandler(h func(appData string) error) {
	c.conn.SetPongHandler(h)
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Connection) SetUserID(uid string) {
	c.uid = uid
}

func (c *Connection) UserID() string {
	return c.uid
}

func (c *Connection) SetRoom(room string) {
	c.room = room
}

func (c *Connection) Room() string {
	return c.room
}

func (c *Connection) SetAppID(appid string) {
	c.appid = appid
}

func (c *Connection) AppID() string {
	return c.appid
}

func (c *Connection) SetProjectID(pid int64) {
	c.pid = pid
}

func (c *Connection) ProjectID() int64 {
	return c.pid
}

func (c *Connection) ID() string {
	return c.id
}

func (c *Connection) Close() error {
	if c.close != nil {
		c.close()
	}
	return c.conn.Close()
}

func (c *Connection) ReadMessage() (int, []byte, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.conn.ReadMessage()
}

func (c *Connection) ReadJSON(v interface{}) error {
	return c.conn.ReadJSON(v)
}

func (c *Connection) WriteMessage(messageType int, data []byte) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.conn.WriteMessage(messageType, data)
}

func (c *Connection) WriteJSON(v interface{}) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.conn.WriteJSON(v)
}

// SetContext ...
func (c *Connection) SetContext(ctx context.Context) {
	if ctx == nil {
		panic("Connection.WithContext must not be nil")
	}
	c.ctx = ctx
}

// Context ...
func (c *Connection) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

type connectionSet struct {
	m map[*Connection]struct{}
}

func newConnectionSet(cap int) connectionSet {
	return connectionSet{
		m: make(map[*Connection]struct{}, cap),
	}
}

func (cs *connectionSet) add(conns ...*Connection) {
	if cs.m == nil {
		cs.m = make(map[*Connection]struct{}, len(conns))
	}
	for _, conn := range conns {
		cs.m[conn] = empty
	}
}

func (cs *connectionSet) each(f func(c *Connection)) {
	for c := range cs.m {
		f(c)
	}
}

func (cs *connectionSet) exist(conn *Connection) bool {
	_, ok := cs.m[conn]
	return ok
}

func (cs *connectionSet) remove(conns ...*Connection) {
	for _, conn := range conns {
		delete(cs.m, conn)
	}
}
