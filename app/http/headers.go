package http

import (
	"strconv"
	"strings"
)

const HeaderContentLength = "content-length"
const HeaderContentType = "content-type"

type Headers struct {
	pairs map[string]string
}

func NewHeaders(pairs map[string]string) *Headers {
	headers := make(map[string]string, len(pairs))
	for k, v := range pairs {
		headers[strings.ToLower(k)] = v
	}
	return &Headers{pairs: headers}
}

func (headers *Headers) Get(key string) (string, bool) {
	value, ok := headers.pairs[strings.ToLower(key)]
	return value, ok
}

func (headers *Headers) Put(key string, value string) {
	headers.pairs[strings.ToLower(key)] = value
}

func (headers *Headers) ContentLength() (int, error) {
	value, ok := headers.Get(HeaderContentLength)
	if !ok {
		return 0, nil
	}
	length, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (headers *Headers) Pairs() map[string]string {
	return headers.pairs
}
