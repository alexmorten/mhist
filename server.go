package mhist

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

//Server is the handler for requests
type Server struct {
	store       *Store
	pools       *Pools
	httpHandler *HTTPHandler
	tcpHandler  *TCPHandler
	waitGroup   *sync.WaitGroup
}

//NewServer returns a new Server
func NewServer(memorySize int) *Server {
	memStore := NewStore(memorySize)
	pools := NewPools(memStore)
	diskStore, err := NewDiskStore(pools)
	if err != nil {
		panic(err)
	}
	memStore.AddSubscriber(diskStore)
	memStore.SetDiskStore(diskStore)

	server := &Server{
		store:     memStore,
		pools:     pools,
		waitGroup: &sync.WaitGroup{},
	}
	tcpHandler := NewTCPHandler(server, 6667)
	server.tcpHandler = tcpHandler
	memStore.AddSubscriber(tcpHandler)

	httpHandler := &HTTPHandler{
		Server: server,
	}
	server.httpHandler = httpHandler

	return server
}

//Run the server
func (s *Server) Run() {
	s.waitGroup.Add(2)
	go func() {
		s.httpHandler.Run()
		s.waitGroup.Done()
	}()
	go func() {
		s.tcpHandler.Run()
		s.waitGroup.Done()
	}()
	s.waitGroup.Wait()
}

//Shutdown all goroutines
func (s *Server) Shutdown() {
	s.store.Shutdown()
}

type message struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (s *Server) handleNewMessage(byteSlice []byte, isReplication bool, onError func(err error, status int)) {
	data := &message{}
	err := json.Unmarshal(byteSlice, data)
	if err != nil {
		onError(err, http.StatusBadRequest)
		return
	}
	if data.Name == "" {
		err = errors.New("name can't be empty")
		onError(err, http.StatusBadRequest)
		return
	}
	measurement, err := s.constructMeasurementFromMessage(data)
	if err != nil {
		onError(err, http.StatusBadRequest)
		return
	}
	s.store.Add(data.Name, measurement, isReplication)
}

func (s *Server) constructMeasurementFromMessage(r *message) (measurement Measurement, err error) {
	switch r.Value.(type) {
	case float64:
		m := s.pools.GetNumericalMeasurement()
		m.Reset()
		m.Ts = time.Now().UnixNano()
		m.Value = r.Value.(float64)
		measurement = m
	case string:
		m := s.pools.GetCategoricalMeasurement()
		m.Reset()
		m.Ts = time.Now().UnixNano()
		m.Value = r.Value.(string)
		measurement = m
	default:
		return nil, errors.New("value is neither a float nor a string")
	}
	return
}
