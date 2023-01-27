package server

import (
	"github.com/ninja-way/pc-store/internal/repository"
	"log"
	"net/http"
)

type Server struct {
	s  *http.Server
	db repository.Repository
}

func New(addr string, handler http.Handler, db repository.Repository) *Server {
	return &Server{
		s: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		db: db,
	}
}

func (s Server) installHandlers() {
	http.HandleFunc("/computers", s.getComputers)
	http.HandleFunc("/computer", s.addComputer)
	http.HandleFunc("/computer/", s.manageComputer)
}

func (s Server) Run() error {
	s.installHandlers()
	log.Println("Server running...")

	if err := s.s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
