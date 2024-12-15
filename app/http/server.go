package http

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Server struct {
	address    string
	dispatcher *Dispatcher
}

func NewServer(address string) *Server {
	return &Server{
		address:    address,
		dispatcher: NewDispatcher(),
	}
}

func (server *Server) AddHandler(uri string, handle func(request *Request, response *Response) error) {
	server.dispatcher.Register(uri, handle)
}

func (server *Server) prepare(conn net.Conn) (*Request, *Response, error) {
	request, err := NewRequest(bufio.NewReader(conn))
	if err != nil {
		return nil, nil, err
	}
	response, err := NewResponse(bufio.NewWriter(conn), request.protocol)
	if err != nil {
		return nil, nil, err
	}

	return request, response, nil
}

func (server *Server) dispatch(conn net.Conn) {
	for {
		request, response, err := server.prepare(conn)
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			conn.Close()
			return
		}
		fmt.Printf("request is: %+v\n", *request)
		fmt.Printf("headers is: %+v\n", *(request.Headers()))
		fmt.Printf("response is: %+v\n", *response)

		handler := server.dispatcher.Dispatch(request)
		fmt.Printf("handler is: %+v\n", handler)
		if handler == nil {
			response.Status404()
			response.End()
			continue
		}
		err = handler.Handle(request, response)
		if err != nil {
			conn.Close()
			fmt.Println("Error accepting connection: ", err)
			return
		}
		connection, ok := response.GetHeader("connection")
		if ok && strings.EqualFold(connection, "close") {
			conn.Close()
			return
		}
	}
}

func (server *Server) Run() {
	l, err := net.Listen("tcp", server.address)
	defer l.Close()
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err)
		}

		go server.dispatch(conn)
	}
}
