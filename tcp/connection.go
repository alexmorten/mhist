package tcp

import (
	"bufio"
	"net"
)

//Connection handles reads and writes to the connection
type Connection struct {
	Socket            net.Conn
	Reader            *bufio.Reader
	onNewMessage      func(message []byte)
	onConnectionClose func()
}

//OnConnectionClose call f when the connection is finished
func (c *Connection) OnConnectionClose(f func()) {
	c.onConnectionClose = f
}

//OnNewMessage call f when the connection is finished
func (c *Connection) OnNewMessage(f func([]byte)) {
	c.onNewMessage = f
}

//Write bytes to connection
func (c *Connection) Write(byteSlice []byte) {
	c.Socket.Write(append(byteSlice, '\n'))
}

//Listen for new messages
func (c *Connection) Listen() {
	for {
		byteSlice, err := c.Reader.ReadSlice('\n')
		if err != nil {
			break
		}
		if c.onNewMessage != nil {
			c.onNewMessage(byteSlice)
		}
	}
	if c.onConnectionClose != nil {
		c.onConnectionClose()
	}
	c.Socket.Close()
}
