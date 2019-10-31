package mhist

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// DebugHandler exposes a debug port over http, used for pprof
type DebugHandler struct {
	Port       int
	httpServer *http.Server
}

// Run listens on the given port and serves http
func (h *DebugHandler) Run() {
	h.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%v", h.Port),
	}

	log.Println("debug_handler running on ", h.httpServer.Addr)
	err := h.httpServer.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

// Shutdown the debug listener
func (h *DebugHandler) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := h.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println("err while shutting debug handler down:", err)
	}
}
