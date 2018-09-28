package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

//Handler for incoming tcp connections
type Handler struct {
	Address            string
	outboundCollection *ConnectionCollection
	onNewMessage       func(message []byte, isReplication bool)
}

//NewHandler with initializaed connection collection, Listen() still needs to be called
func NewHandler(address string) *Handler {
	return &Handler{
		Address:            address,
		outboundCollection: &ConnectionCollection{},
	}
}

//OnNewMessage calls f with the message and states wether or not the message is from replication
func (h *Handler) OnNewMessage(f func(message []byte, isReplication bool)) {
	h.onNewMessage = f
}

//Notify handler about new message
func (h *Handler) Notify(message []byte) {
	h.outboundCollection.ForEach(func(conn *Connection) {
		conn.Write(message)
	})
}

//Listen for new connections
func (h *Handler) Listen() {
	listener, err := net.Listen("tcp", h.Address)
	if err != nil {
		panic("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		reader := bufio.NewReader(conn)
		byteSlice, err := reader.ReadSlice('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		m := &SubscribtionMessage{}
		err = json.Unmarshal(byteSlice, m)
		if err != nil {
			fmt.Println(err)
			continue
		}
		connectionWrapper := &Connection{
			Socket: conn,
		}
		if m.Publisher {
			connectionWrapper.OnNewMessage(func(byteSlice []byte) {
				if h.onNewMessage != nil {
					h.onNewMessage(byteSlice, m.Replication)
				}
			})
		} else {
			h.outboundCollection.AddConnection(connectionWrapper)
			connectionWrapper.OnConnectionClose(func() {
				h.outboundCollection.RemoveConnection(connectionWrapper)
			})
		}
	}
}
