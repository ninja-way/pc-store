package transport

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Server with custom gin router
type Server struct {
	s *http.Server
}

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
			log.WithField("listen server", err).Fatal()
		}
	}()
	log.Warn("Server running...")

	<-ctx.Done()
	cancel()

	// gracefully shutdown the server
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.s.Shutdown(ctx); err != nil {
		log.WithField("server shutdown", err).Warn()
	}

	log.Warn("Server stopped.")
}
