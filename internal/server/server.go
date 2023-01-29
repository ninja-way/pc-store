package server

import (
	"github.com/ninja-way/pc-store/internal/repository"
	"log"
	"net/http"
)

type Server struct {
	s  *http.Server
	db repository.DB
}

func New(addr string, handler http.Handler, db repository.DB) *Server {
	return &Server{
		s: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		db: db,
	}
}

func (s Server) setupHandlers() {
	http.HandleFunc("/computers", s.getComputers)
	http.HandleFunc("/computer", s.addComputer)
	http.HandleFunc("/computer/", s.manageComputer)
}

func (s Server) Run() error {
	s.setupHandlers()
	log.Println("Server running...")

	if err := s.s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
