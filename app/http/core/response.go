package core

import (
	"bufio"
	"fmt"
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"strconv"
	"strings"
)

type Response struct {
	writer   *bufio.Writer
	protocol string
	status   uint32
	message  string
	headers  *Headers
	body     []byte
}

func (response *Response) Protocol() string {
	return response.protocol
}

func (response *Response) Status() uint32 {
	return response.status
}

func (response *Response) SetStatus(status uint32) {
	response.status = status
}

func (response *Response) Message() string {
	return response.message
}

func (response *Response) SetMessage(message string) {
	response.message = message
}

func (response *Response) GetHeader(key string) (string, bool) {
	return response.headers.Get(key)
}

func (response *Response) SetHeader(key string, value string) {
	response.headers.Put(strings.ToLower(key), value)
}

func (response *Response) SetBody(contentType string, encoding string, body []byte) {
	if encoding != "" {
		contentType += "; charset=" + encoding
	}
	response.SetHeader(common.HeaderContentType, contentType)
	if body == nil {
		return
	}
	response.SetHeader(common.HeaderContentLength, strconv.Itoa(len(body)))
	response.body = body
}

func (response *Response) End() {
	if response.body == nil {
		response.headers.Put(common.HeaderContentLength, "0")
	}
	responseLine := fmt.Sprintf("%s %d %s", response.protocol, response.status, response.message)
	response.writer.WriteString(responseLine)
	response.endLine()
	for key, value := range response.headers.Pairs() {
		response.writer.WriteString(fmt.Sprintf("%s: %s", strings.ToLower(key), value))
		response.endLine()
	}
	response.endLine()
	if len(response.body) > 0 {
		response.writer.Write(response.body)
	}

	_ = response.writer.Flush()
}

func (response *Response) endLine() {
	response.writer.WriteString("\r\n")
}

func (response *Response) String() string {
	return fmt.Sprintf("{protocol=%s, status=%d, message=%s, headers=%v, bodyLength=%d}",
		response.protocol, response.status, response.message, response.headers, len(response.body))
}

func NewResponse(writer *bufio.Writer, protocol string) (*Response, error) {
	headers := NewHeaders(make(map[string]string))
	return &Response{
		writer:   writer,
		protocol: protocol,
		headers:  headers,
	}, nil
}
