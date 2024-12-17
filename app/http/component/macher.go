package component

import (
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"regexp"
)

type Matcher interface {
	Match(request *core.Request) bool
}

func matchMethod(method string, methods []string) bool {
	if len(methods) == 0 {
		return true
	}

	for _, m := range methods {
		if m == method {
			return true
		}
	}

	return false
}

type SimpleMatcher struct {
	uri     string
	methods []string
}

func (matcher *SimpleMatcher) Match(request *core.Request) bool {
	if matcher.uri != request.Uri() {
		return false
	}

	return matchMethod(request.Method(), matcher.methods)
}
func NewSimpleMatcher(uri string, methods []string) *SimpleMatcher {
	return &SimpleMatcher{
		uri:     uri,
		methods: methods,
	}
}

type PatternMatcher struct {
	uriPattern *regexp.Regexp
	methods    []string
}

func (matcher *PatternMatcher) Match(request *core.Request) bool {
	if !matcher.uriPattern.MatchString(request.Uri()) {
		return false
	}

	return matchMethod(request.Method(), matcher.methods)
}

func NewPatternMatcher(pattern string, methods []string) (*PatternMatcher, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &PatternMatcher{
		uriPattern: reg,
		methods:    methods,
	}, nil
}
