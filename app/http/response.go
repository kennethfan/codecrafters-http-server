package http

import (
	"bufio"
	"fmt"
	"strings"
)

type Response struct {
	writer   *bufio.Writer
	protocol string
	status   uint32
	message  string
	headers  *Headers
	body     interface{}
}

func NewResponse(writer *bufio.Writer, protocol string) (*Response, error) {
	headers := NewHeaders(make(map[string]string))
	return &Response{
		writer:   writer,
		protocol: protocol,
		headers:  headers,
	}, nil
}

func (response *Response) StatusOK() {
	response.status = 200
	response.message = "OK"
	fmt.Sprintf("200 response is %+v\n", *response)
	response.End()
}

func (response *Response) Status404() {
	response.status = 404
	response.message = "Not Found"
	fmt.Sprintf("404 response is %+v\n", *response)
	response.End()
}

func (response *Response) End() {
	fmt.Sprintf("end response is %+v\n", *response)
	responseLine := fmt.Sprintf("%s %d %s", response.protocol, response.status, response.message)
	fmt.Printf("Response line : %s\n", responseLine)
	response.writer.WriteString(responseLine)
	response.endLine()
	for key, value := range response.headers.Pairs() {
		response.writer.WriteString(fmt.Sprintf("%s: %s", strings.ToLower(key), value))
		response.endLine()
	}
	response.endLine()
	_ = response.writer.Flush()
}

func (response *Response) endLine() {
	response.writer.WriteString("\r\n")
}
