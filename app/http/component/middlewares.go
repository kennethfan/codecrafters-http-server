package component

import (
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"log"
	"sort"
)

type MiddlewareFunc func(request *core.Request, response *core.Response, chain MiddlewareChain) (bool, error)

type MiddlewareChain interface {
	DoChain(request *core.Request, response *core.Response) (bool, error)
}

type Middleware interface {
	DoChain(request *core.Request, response *core.Response, chain MiddlewareChain) (bool, error)
}

type middlewareImpl struct {
	middlewareFunc MiddlewareFunc
}

func (middlewareImpl *middlewareImpl) DoChain(request *core.Request, response *core.Response, chain MiddlewareChain) (bool, error) {
	return middlewareImpl.middlewareFunc(request, response, chain)
}

func NewMiddleware(middlewareFunc MiddlewareFunc) Middleware {
	middleware := &middlewareImpl{
		middlewareFunc: middlewareFunc,
	}

	return middleware
}

type MatchMiddleware struct {
	priority   int
	matcher    Matcher
	middleware Middleware
}

type MatchMiddlewareList []*MatchMiddleware

func (matchMiddlewareList MatchMiddlewareList) Len() int {
	return len(matchMiddlewareList)
}

func (matchMiddlewareList MatchMiddlewareList) Less(i, j int) bool {
	return matchMiddlewareList[i].priority < matchMiddlewareList[j].priority
}

func (matchMiddlewareList MatchMiddlewareList) Swap(i, j int) {
	matchMiddlewareList[i], matchMiddlewareList[j] = matchMiddlewareList[j], matchMiddlewareList[i]
}

func NewMatchMiddleware(priority int, matcher Matcher, middleware Middleware) *MatchMiddleware {
	return &MatchMiddleware{
		priority:   priority,
		matcher:    matcher,
		middleware: middleware,
	}
}

type MiddlewareChainImpl struct {
	matchMiddlewares MatchMiddlewareList
	pos              int
	handler          Handler
}

func (chain MiddlewareChainImpl) DoChain(request *core.Request, response *core.Response) (bool, error) {
	for {
		if chain.pos >= len(chain.matchMiddlewares) {
			break
		}
		matchMiddleware := chain.matchMiddlewares[chain.pos]
		chain.pos += 1
		if matchMiddleware.matcher.Match(request) {
			return matchMiddleware.middleware.DoChain(request, response, chain)
		}
	}

	err := chain.handler.Handle(request, response)
	if err != nil {
		return false, err
	}

	return true, nil
}

type Middlewares struct {
	matchMiddlewareList MatchMiddlewareList
}

func (middlewares *Middlewares) addMiddleware(priority int, matcher Matcher, middleware Middleware) {
	middlewares.matchMiddlewareList = append(middlewares.matchMiddlewareList, NewMatchMiddleware(priority, matcher, middleware))
	sort.Sort(middlewares.matchMiddlewareList)
}

func (middlewares *Middlewares) RegisterSimpleMiddleware(priority int, uri string, middleware Middleware, methods ...string) {
	middlewares.addMiddleware(priority, NewSimpleMatcher(uri, methods), middleware)
}

func (middlewares *Middlewares) RegisterSimpleMiddlewareFunc(priority int, uri string, middlewareFunc MiddlewareFunc, methods ...string) {
	middlewares.RegisterSimpleMiddleware(priority, uri, NewMiddleware(middlewareFunc), methods...)
}

func (middlewares *Middlewares) RegisterPatternMiddleware(priority int, pattern string, middleware Middleware, methods ...string) {
	matcher, err := NewPatternMatcher(pattern, methods)
	if err != nil {
		log.Fatal(err)
	}

	middlewares.addMiddleware(priority, matcher, middleware)
}

func (middlewares *Middlewares) RegisterPatternMiddlewareFunc(priority int, pattern string, middlewareFunc MiddlewareFunc, methods ...string) {
	middlewares.RegisterPatternMiddleware(priority, pattern, NewMiddleware(middlewareFunc), methods...)
}

func (middlewares *Middlewares) NewMiddlewareChain(handler Handler) MiddlewareChain {
	return &MiddlewareChainImpl{
		matchMiddlewares: middlewares.matchMiddlewareList,
		pos:              0,
		handler:          handler,
	}
}

func NewMiddlewares() *Middlewares {
	return &Middlewares{
		matchMiddlewareList: []*MatchMiddleware{},
	}
}
