package server

import (
	"context"
	"github.com/ninja-way/pc-store/internal/repository"
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
	db     repository.DB
	errors chan error
}

// New is server constructor
func New(addr string, handler http.Handler, db repository.DB) *Server {
	return &Server{
		s: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		db:     db,
		errors: make(chan error),
	}
}

func (s Server) setupHandlers() {
	http.HandleFunc("/computers", s.getComputers)
	http.HandleFunc("/computer", s.addComputer)
	http.HandleFunc("/computer/", s.manageComputer)
}

// Run configures and starts the server
func (s Server) Run() error {
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
func (s Server) stop() {
	log.Println("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.s.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}

	if err := s.db.Close(ctx); err != nil {
		log.Printf("Database closing: %v", err)
	}

	close(s.errors)
}
