package mhist

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//Server is the handler for requests
type Server struct {
	store        *Store
	grpcHandler  *GrpcHandler
	debugHandler *DebugHandler
	waitGroup    *sync.WaitGroup
}

//ServerConfig ...
type ServerConfig struct {
	GrpcPort   int
	DebugPort  int
	MemorySize int
	DiskSize   int
}

//NewServer returns a new Server
func NewServer(config ServerConfig) *Server {
	diskStore, err := NewDiskStore(config.MemorySize, config.DiskSize)
	if err != nil {
		panic(err)
	}

	store := NewStore(diskStore)

	server := &Server{
		store:     store,
		waitGroup: &sync.WaitGroup{},
	}

	grpcHandler := NewGrpcHandler(server, config.GrpcPort)
	server.grpcHandler = grpcHandler
	store.AddSubscriber(grpcHandler)

	server.debugHandler = &DebugHandler{
		Port:   config.DebugPort,
		server: server,
	}

	return server
}

//Run the server
func (s *Server) Run() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		signal := <-signals
		log.Printf("received %s, shutting down\n", signal)
		s.Shutdown()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		s.grpcHandler.Run()
		wg.Done()
	}()
	go func() {
		s.debugHandler.Run()
		wg.Done()
	}()

	wg.Wait()
}

//Shutdown all goroutines and connections
func (s *Server) Shutdown() {
	s.grpcHandler.Shutdown()
	s.debugHandler.Shutdown()

	s.store.Shutdown()
}
