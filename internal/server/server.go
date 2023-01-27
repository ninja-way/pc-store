package server

import "net/http"

type Server struct {
	s *http.Server
}

func New() *Server {
	return &Server{
		&http.Server{},
	}
}
