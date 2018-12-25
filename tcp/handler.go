package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/alexmorten/mhist/models"
)

//MessageHandler handles everything for incoming messages
type MessageHandler interface {
	HandleNewMessage(byteSlice []byte, isReplication bool, onError func(err error, _ int))
}

//Handler handles tcp connections
type Handler struct {
	messageHandler              MessageHandler
	address                     string
	outboundCollection          *ConnectionCollection
	filterPerOutboundConnection map[*Connection]*models.FilterCollection
	filterMutex                 *sync.RWMutex
	pools                       *models.Pools
}

//NewHandler sets the wrapped handlers callbacks correctly, Run() still has to be called
func NewHandler(port int, messageHandler MessageHandler, pools *models.Pools) *Handler {
	return &Handler{
		messageHandler:              messageHandler,
		address:                     fmt.Sprintf("0.0.0.0:%v", port),
		outboundCollection:          &ConnectionCollection{},
		filterMutex:                 &sync.RWMutex{},
		filterPerOutboundConnection: make(map[*Connection]*models.FilterCollection),
		pools:                       pools,
	}
}

//Notify handler about new message
func (h *Handler) Notify(name string, measurement models.Measurement) {
	m := h.pools.GetMessage()
	defer h.pools.PutMessage(m)

	m.Reset()
	m.Name = name
	m.Value = measurement.ValueInterface()
	m.Timestamp = measurement.Timestamp()

	byteSlice, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return
	}
	h.filterMutex.RLock()
	defer h.filterMutex.RUnlock()
	h.outboundCollection.ForEach(func(conn *Connection) {
		filter := h.filterPerOutboundConnection[conn]
		if filter != nil {
			if filter.Passes(name, measurement) {
				conn.Write(byteSlice)
			}
		} else {
			fmt.Println("Filter for outbound connection was nil, please investigate!")
		}
	})
}

//Run listens for new connections
func (h *Handler) Run() {
	listener, err := net.Listen("tcp", h.address)
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
		go h.handleNewConnection(conn)
	}
}

func (h *Handler) onNewMessage(byteSlice []byte, isReplication bool) {
	h.messageHandler.HandleNewMessage(byteSlice, isReplication, func(err error, _ int) {
		if err != nil {
			fmt.Println(err)
		}
	})
}

func (h *Handler) handleNewConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	byteSlice, err := reader.ReadSlice('\n')
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}
	m := &models.SubscriptionMessage{}
	err = json.Unmarshal(byteSlice, m)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}
	connectionWrapper := &Connection{
		Socket: conn,
		Reader: reader,
	}
	if m.Publisher {
		connectionWrapper.OnNewMessage(func(byteSlice []byte) {
			h.onNewMessage(byteSlice, m.Replication)
		})
	} else {
		h.addFilterForConnection(m.FilterDefinition, connectionWrapper)
		h.outboundCollection.AddConnection(connectionWrapper)
		connectionWrapper.OnConnectionClose(func() {
			h.outboundCollection.RemoveConnection(connectionWrapper)
			h.removeFilterForConnection(connectionWrapper)
		})
	}
	connectionWrapper.Listen()
}

func (h *Handler) removeFilterForConnection(conn *Connection) {
	h.filterMutex.Lock()
	defer h.filterMutex.Unlock()

	delete(h.filterPerOutboundConnection, conn)
}

func (h *Handler) addFilterForConnection(filterDefinition models.FilterDefinition, conn *Connection) {
	h.filterMutex.Lock()
	defer h.filterMutex.Unlock()

	filter := models.NewFilterCollection(filterDefinition)
	h.filterPerOutboundConnection[conn] = filter
}
