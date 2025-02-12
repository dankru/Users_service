package server

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	server http.Server
}

func NewServer(addr string, writeTimeout, readTimeout, idleTimeout time.Duration, handler http.Handler) Server {
	server := Server{
		server: http.Server{
			Addr:         addr,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			Handler:      handler,
		},
	}
	return server
}

func (s *Server) Run() {
	err := s.server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to run server", err.Error())
	}
}
