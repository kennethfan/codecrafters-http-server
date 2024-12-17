package core

import (
	"bufio"
	"fmt"
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"log"
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

func NewResponse(writer *bufio.Writer, protocol string) (*Response, error) {
	headers := NewHeaders(make(map[string]string))
	return &Response{
		writer:   writer,
		protocol: protocol,
		headers:  headers,
	}, nil
}

func (response *Response) StatusOK() {
	response.status = common.StatusOK
	response.message = common.MessageOk
	log.Printf("200 response is %+v\n", *response)
}

func (response *Response) Status404() {
	response.status = common.StatusNotFound
	response.message = common.MessageNotFound
	log.Printf("404 response is %+v\n", *response)
}

func (response *Response) GetHeader(key string) (string, bool) {
	return response.headers.Get(key)
}

func (response *Response) SetHeader(key string, value string) {
	response.headers.Put(strings.ToLower(key), value)
}

func (response *Response) Html(html string, encoding string) {
	response.setBody("text/html", encoding, []byte(html))
}
func (response *Response) Json(json string, encoding string) {
	response.setBody("application/json", encoding, []byte(json))
}

func (response *Response) PlainText(text string, encoding string) {
	response.setBody("text/plain;", encoding, []byte(text))
}

func (response *Response) setBody(contentType string, encoding string, body []byte) {
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
	log.Printf("end response is %+v\n", *response)
	responseLine := fmt.Sprintf("%s %d %s", response.protocol, response.status, response.message)
	log.Printf("Response line : %s\n", responseLine)
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
