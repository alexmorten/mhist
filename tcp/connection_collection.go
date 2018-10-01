package tcp

import (
	"sync"
)

//ConnectionCollection handles all active connections for outgoing messages
type ConnectionCollection struct {
	connections []*Connection
	sync.RWMutex
}

//AddConnection to collection
func (c *ConnectionCollection) AddConnection(addConnection *Connection) {
	c.Lock()
	defer c.Unlock()

	for index, connection := range c.connections {
		if connection == nil {
			c.connections[index] = addConnection
			return
		}
	}
	c.connections = append(c.connections, addConnection)
}

//RemoveConnection from collection
func (c *ConnectionCollection) RemoveConnection(removeConnection *Connection) {
	c.Lock()
	defer c.Unlock()

	for index, connection := range c.connections {
		if connection == removeConnection {
			c.connections[index] = nil
			return
		}
	}
}

//ForEach connection in the collection
func (c *ConnectionCollection) ForEach(f func(conn *Connection)) {
	c.RLock()
	defer c.RUnlock()

	for _, conn := range c.connections {
		if conn != nil {
			f(conn)
		}
	}
}
