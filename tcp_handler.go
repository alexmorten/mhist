package mhist

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/codeuniversity/ppp-mhist/tcp"
)

//TCPHandler handles tcp connections. This a wrapper around tcp.handler with (Un-)Marshalizing capabilities
type TCPHandler struct {
	server         *Server
	wrappedHandler *tcp.Handler
}

//NewTCPHandler sets the wrapped handlers callbacks correctly, Run() still has to be called
func NewTCPHandler(server *Server, port int) *TCPHandler {
	wrappedHandler := tcp.NewHandler("localhost:" + strconv.FormatInt(int64(port), 10))
	wrappedHandler.OnNewMessage(func(byteSlice []byte, isReplication bool) {
		server.handleNewMessage(byteSlice, isReplication, func(err error, _ int) {
			if err != nil {
				fmt.Println(err)
			}
		})
	})
	return &TCPHandler{
		server:         server,
		wrappedHandler: wrappedHandler,
	}
}

//Run the TCPHandler
func (h *TCPHandler) Run() {
	h.wrappedHandler.Listen()
}

//Notify handler about new message
func (h *TCPHandler) Notify(name string, measurement Measurement) {
	m := &message{
		Name:  name,
		Value: measurement.ValueInterface(),
	}

	byteSlice, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	h.wrappedHandler.Notify(byteSlice)
}
