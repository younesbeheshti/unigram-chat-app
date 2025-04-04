package config

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

	fmt.Println("server is up on -> 127.0.0.1:15000")

	log.Fatal(http.ListenAndServe(s.Port, router))
}