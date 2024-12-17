package core

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	reader     *bufio.Reader
	method     string
	url        string
	protocol   string
	headers    *Headers
	bodyLength int
	body       interface{}
}

func (request *Request) String() string {
	return fmt.Sprintf("{protocol=%s, method=%s, url=%s, headers=%v, bodyLength=%d}",
		request.protocol, request.method, request.url, request.headers, request.bodyLength)
}

func NewRequest(reader *bufio.Reader) (*Request, error) {
	requestLine, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	parts := strings.Fields(requestLine)
	if len(parts) != 3 {
		return nil, errors.New(requestLine + " not a valid request line")
	}
	method := strings.ToUpper(parts[0])
	url := parts[1]
	protocol := parts[2]
	headerPairs, err := readHeaders(reader)
	if err != nil {
		return nil, err
	}
	headers := NewHeaders(headerPairs)
	contentLength, err := headers.ContentLength()
	if err != nil {
		return nil, err
	}

	return &Request{
		reader:     reader,
		method:     method,
		url:        url,
		protocol:   protocol,
		headers:    headers,
		bodyLength: contentLength,
		body:       nil,
	}, nil
}

func (request *Request) Uri() string {
	parts := strings.SplitN(request.url, "?", 2)
	return parts[0]
}

func (request *Request) Method() string {
	return request.method
}

func (request *Request) Protocol() string {
	return request.protocol
}

func (request *Request) Headers() *Headers {
	return request.headers
}

func (request *Request) GetHeader(key string) (string, bool) {
	return request.headers.Get(key)
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func readHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			panic(line + " not a header line")
		}
		headers[strings.ToLower(parts[0])] = parts[1]
	}

	return headers, nil
}
