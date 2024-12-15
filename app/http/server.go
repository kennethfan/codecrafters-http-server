package http

import (
	"bufio"
	"fmt"
	"net"
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
	request, response, err := server.prepare(conn)
	fmt.Printf("request is: %+v\n", *request)
	fmt.Printf("headers is: %+v\n", *(request.Headers()))
	fmt.Printf("response is: %+v\n", *response)
	if err != nil {
		fmt.Println("Error accepting connection: ", err)
		return
	}

	handler := server.dispatcher.Dispatch(request)
	fmt.Printf("handler is: %+v\n", handler)
	if handler == nil {
		response.Status404()
		return
	}
	err = handler.Handle(request, response)
	if err != nil {
		fmt.Println("Error accepting connection: ", err)
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
