package transport

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server stores the config and connected database
type Server struct {
	s      *http.Server
	h      *Handler
	errors chan error
}

// NewServer is server constructor
func NewServer(addr string, mw http.Handler, h *Handler) *Server {
	return &Server{
		s: &http.Server{
			Addr:    addr,
			Handler: mw,
		},
		h:      h,
		errors: make(chan error),
	}
}

func (s *Server) setupHandlers() {
	http.HandleFunc("/computers", s.h.getComputers)
	http.HandleFunc("/computer", s.h.addComputer)
	http.HandleFunc("/computer/", s.h.manageComputer)
}

// Run configures and starts the server
func (s *Server) Run() error {
	s.setupHandlers()

	go func() {
		if err := s.s.ListenAndServe(); err != nil {
			s.errors <- err
		}
	}()

	// chanel for track interrupts
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	log.Println("Server running...")

	// stop server in both cases
	select {
	case <-quit:
		s.stop()
	case err := <-s.errors:
		s.stop()
		return err
	}
	return nil
}

// Shuts down the server and close db connection
func (s *Server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.s.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}

	//close(s.errors)
	log.Println("Server stopped")
}
