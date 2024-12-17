package http

import (
	"bufio"
	"github.com/kennethfan/codecrafters-http-server/http/component"
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"log"
	"net"
	"strings"
)

type Worker struct {
	router      component.Router
	middlewares *component.Middlewares
}

func (worker *Worker) prepare(conn net.Conn) (*core.Request, *core.Response, error) {
	request, err := core.NewRequest(bufio.NewReader(conn))
	if err != nil {
		return nil, nil, err
	}
	response, err := core.NewResponse(bufio.NewWriter(conn), request.Protocol())
	if err != nil {
		return nil, nil, err
	}

	return request, response, nil
}

func (worker *Worker) work(conn net.Conn) {
	for {
		request, response, err := worker.prepare(conn)
		if err != nil {
			log.Println("Error accepting connection: ", err)
			conn.Close()
			return
		}
		log.Printf("request is: %+v\n", *request)
		log.Printf("headers is: %+v\n", *(request.Headers()))
		log.Printf("response is: %+v\n", *response)

		handler := worker.router.Dispatch(request)
		log.Printf("handler is: %+v\n", handler)
		if handler == nil {
			response.Status404()
			response.End()
			continue
		}

		chain := worker.middlewares.NewMiddlewareChain(handler)
		ok, err := chain.DoChain(request, response)
		if err != nil {
			conn.Close()
			log.Println("Error accepting connection: ", err)
			return
		}
		connection, ok := response.GetHeader("connection")
		if ok && strings.EqualFold(connection, "close") {
			conn.Close()
			return
		}
	}
}

func NewWorker(router component.Router, middlewares *component.Middlewares) *Worker {
	return &Worker{
		router:      router,
		middlewares: middlewares,
	}
}
