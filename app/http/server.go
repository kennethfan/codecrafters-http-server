package http

import (
	"github.com/kennethfan/codecrafters-http-server/http/component"
	"log"
	"net"
)

type Server struct {
	address     string
	router      component.Router
	middlewares *component.Middlewares
}

type ServerOption func(*Server)

func WithAddress(address string) ServerOption {
	return func(server *Server) {
		server.address = address
	}
}

func WithRouter(router component.Router) ServerOption {
	return func(server *Server) {
		server.router = router
	}
}

func NewServer(options ...ServerOption) *Server {
	server := new(Server)
	for _, option := range options {
		option(server)
	}
	if server.address == "" {
		server.address = "0.0.0.0:4221"
	}
	if server.router == nil {
		server.router = component.NewRouter()
	}
	server.middlewares = component.NewMiddlewares()
	return server
}

func (server *Server) Router() component.Router {
	return server.router
}

func (server *Server) Middlewares() *component.Middlewares {
	return server.middlewares
}

func (server *Server) work(conn net.Conn) {
	NewWorker(server.router, server.middlewares).work(conn)
}

func (server *Server) Start() {
	l, err := net.Listen("tcp", server.address)
	defer l.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Println("Error accepting connection: ", err)
		}

		go server.work(conn)
	}
}
