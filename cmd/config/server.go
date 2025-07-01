package config

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/younesbeheshti/chatapp-backend/cmd/middleware"
	"log"
	"net/http"
)

type Server struct {
	Address string
	Port    string
	Type    string
}

// make a new server
func NewServer(Address, Port, Type string) *Server {
	return &Server{
		Address: Address,
		Port:    Port,
		Type:    Type,
	}
}

// function to init(or run) the server
func (s *Server) InitServer(router *mux.Router) {

	fmt.Printf("server is up on -> 127.0.0.1:%s", s.Port)

	log.Fatal(http.ListenAndServe(s.Port, middleware.CorsHandler(router)))
}
