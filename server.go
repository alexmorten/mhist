package mhist

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/alexmorten/mhist/models"
	"github.com/alexmorten/mhist/tcp"
)

//Server is the handler for requests
type Server struct {
	store       *Store
	pools       *models.Pools
	httpHandler *HTTPHandler
	tcpHandler  *tcp.Handler
	waitGroup   *sync.WaitGroup
}

//ServerConfig ...
type ServerConfig struct {
	HTTPPort             int
	TCPPort              int
	MemorySize           int
	DiskSize             int
	ReplicationAddresses []string
}

//NewServer returns a new Server
func NewServer(config ServerConfig) *Server {
	pools := models.NewPools()
	diskStore, err := NewDiskStore(pools, config.MemorySize, config.DiskSize)
	if err != nil {
		panic(err)
	}

	store := NewStore(diskStore)

	server := &Server{
		store:     store,
		pools:     pools,
		waitGroup: &sync.WaitGroup{},
	}
	tcpHandler := tcp.NewHandler(config.TCPPort, server, pools)
	server.tcpHandler = tcpHandler
	store.AddSubscriber(tcpHandler)

	httpHandler := &HTTPHandler{
		Server: server,
		Port:   config.HTTPPort,
	}
	server.httpHandler = httpHandler
	for _, address := range config.ReplicationAddresses {
		replication := tcp.NewReplication(address, pools)
		store.AddReplication(replication)
	}
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
}

//HandleNewMessage coming from any source
func (s *Server) HandleNewMessage(byteSlice []byte, isReplication bool, onError func(err error, status int)) {
	data := s.pools.GetMessage()
	defer s.pools.PutMessage(data)

	data.Reset()
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

func (s *Server) constructMeasurementFromMessage(message *models.Message) (measurement models.Measurement, err error) {
	switch message.Value.(type) {
	case float64:
		m := s.pools.GetNumericalMeasurement()
		m.Reset()
		if message.Timestamp == 0 {
			m.Ts = time.Now().UnixNano()
		} else {
			m.Ts = message.Timestamp
		}
		m.Value = message.Value.(float64)
		measurement = m
	case string:
		m := s.pools.GetCategoricalMeasurement()
		m.Reset()
		if message.Timestamp == 0 {
			m.Ts = time.Now().UnixNano()
		} else {
			m.Ts = message.Timestamp
		}
		m.Value = message.Value.(string)
		measurement = m
	default:
		return nil, errors.New("value is neither a float nor a string")
	}
	return
}
