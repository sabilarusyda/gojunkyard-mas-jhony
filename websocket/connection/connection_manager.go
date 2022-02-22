package connection

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Manager ...
// Case:
// Connection 1 		=> Channel Video A Notification, Channel Video A Chat, Channel Global Notification
// Connection 2 		=> Channel Global Notification
// Channel Video A 		=> Connection 1, Connection 2
// Solution:
// 1. Connection holds list of channel pointer
// 2. Channel holds list of connection
type Manager struct {
	mux                sync.RWMutex
	conns              connectionSet
	channels           map[string]*channel
	connectionUpgrader *websocket.Upgrader
}

func NewManager() *Manager {
	return &Manager{
		conns:    newConnectionSet(1000),
		channels: make(map[string]*channel, 1000),
	}
}

func (cm *Manager) UpgradeConnection(w http.ResponseWriter, r *http.Request, option *Option) (*Connection, error) {
	conn, err := New(w, r, option)
	if err != nil {
		return nil, err
	}
	conn.close = func() { cm.UnregisterConnection(conn) }

	cm.mux.Lock()
	cm.conns.add(conn)
	cm.mux.Unlock()
	return conn, err
}

func (cm *Manager) RegisterChannel(conn *Connection, channel string) {
	// step 1. Ignore if connection is not registered to Connection Manager
	if !cm.conns.exist(conn) {
		return
	}

	// step 2. Create New Channel if not exist
	cm.mux.RLock()
	ch, ok := cm.channels[channel]
	cm.mux.RUnlock()
	if !ok {
		ch = newChannel()

		cm.mux.Lock()
		cm.channels[channel] = ch
		cm.mux.Unlock()
	}

	// step 3. Register connection to Channel Subscriber List
	ch.conns.add(conn)

	// step 4. Register channel to Connection Channel List
	conn.channels.add(ch)
}

func (cm *Manager) BroadcastMessage(channel string, messageType int, data []byte) {
	// step 1. Get Channel and return if channel not exist
	cm.mux.Lock()
	ch, ok := cm.channels[channel]
	cm.mux.Unlock()
	if !ok {
		return
	}

	// step 2. Broadcast data to all connection in channel
	ch.conns.each(func(c *Connection) {
		c.conn.WriteMessage(messageType, data)
	})
}

func (cm *Manager) BroadcastJSON(channel string, v interface{}) {
	// step 1. Get Channel and return if channel not exist
	cm.mux.Lock()
	ch, ok := cm.channels[channel]
	cm.mux.Unlock()
	if !ok {
		return
	}

	// step 2. Broadcast data to all connection in channel
	ch.conns.each(func(c *Connection) {
		time.Sleep(100 * time.Nanosecond)
		c.conn.WriteJSON(v)
	})
}

func (cm *Manager) UnregisterChannel(conn *Connection, channel string) {
	// step 1. Remove connection from Connection Manager
	cm.conns.remove(conn)

	// step 2. Remove connection from Channel Subscriber List
	cm.mux.RLock()
	ch, ok := cm.channels[channel]
	cm.mux.RUnlock()
	if ok {
		ch.conns.remove(conn)
	}
}

func (cm *Manager) UnregisterConnection(conn *Connection) {
	cm.conns.remove(conn)
	conn.channels.each(func(c *channel) {
		c.conns.remove(conn)
	})
}
