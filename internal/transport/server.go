package transport

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Server stores the config and connected database
type Server struct {
	s *http.Server
}

// NewServer is server constructor
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		s: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// Run configures and starts the server
func (s *Server) Run() {
	// track interrupts
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// run server
	go func() {
		if err := s.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	log.Println("Server running...")

	<-ctx.Done()
	cancel()

	// gracefully shutdown the server
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.s.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}

	log.Println("Server stopped")
}
