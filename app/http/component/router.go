package component

import (
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"log"
	"sort"
)

type HandleFunc func(*core.Request, *core.Response) error

type Handler interface {
	Handle(*core.Request, *core.Response) error
}

type handlerImpl struct {
	handleFunc HandleFunc
}

func (handler *handlerImpl) Handle(request *core.Request, response *core.Response) error {
	err := handler.handleFunc(request, response)
	return err
}

func NewHandler(handleFunc HandleFunc) Handler {
	return &handlerImpl{handleFunc: handleFunc}
}

type Router interface {
	RegisterSimpleHandler(uri string, handler Handler, methods ...string)
	RegisterSimpleHandleFunc(uri string, handleFunc HandleFunc, methods ...string)
	RegisterPatternHandler(priority int, pattern string, handler Handler, methods ...string)
	RegisterPatternHandleFunc(priority int, pattern string, handleFunc HandleFunc, methods ...string)
	Dispatch(request *core.Request) Handler
}

type PatternHandler struct {
	priority       int
	patternMatcher *PatternMatcher
	handler        Handler
}

type PatternHandlerList []*PatternHandler

func (patternHandlerList PatternHandlerList) Len() int {
	return len(patternHandlerList)
}

func (patternHandlerList PatternHandlerList) Less(i, j int) bool {
	return patternHandlerList[i].priority < patternHandlerList[j].priority
}

func (patternHandlerList PatternHandlerList) Swap(i, j int) {
	patternHandlerList[i], patternHandlerList[j] = patternHandlerList[j], patternHandlerList[i]
}

func newPatternHandler(priority int, patternMatcher *PatternMatcher, handler Handler) *PatternHandler {
	return &PatternHandler{
		priority:       priority,
		patternMatcher: patternMatcher,
		handler:        handler,
	}
}

type defaultRouter struct {
	simpleHandlers  map[string]map[string]Handler
	patternHandlers PatternHandlerList
}

func (router *defaultRouter) RegisterSimpleHandler(uri string, handler Handler, methods ...string) {
	if len(methods) == 0 {
		methods = common.HttpCommonMethods
	}

	methodHandlers, ok := router.simpleHandlers[uri]
	if !ok {
		methodHandlers = make(map[string]Handler)
		for _, method := range methods {
			methodHandlers[method] = handler
		}
		router.simpleHandlers[uri] = methodHandlers
		return
	}

	for _, method := range methods {
		_, ok := methodHandlers[method]
		if ok {
			log.Printf("[%s] %s duplicate", method, uri)
		}

		methodHandlers[method] = handler
	}
}

func (router *defaultRouter) RegisterSimpleHandleFunc(uri string, handleFunc HandleFunc, methods ...string) {
	router.RegisterSimpleHandler(uri, NewHandler(handleFunc), methods...)
}

func (router *defaultRouter) RegisterPatternHandler(priority int, pattern string, handler Handler, methods ...string) {
	patternMatcher, err := NewPatternMatcher(pattern, methods)
	if err != nil {
		log.Fatal(err)
	}

	router.patternHandlers = append(router.patternHandlers, newPatternHandler(priority, patternMatcher, handler))
	sort.Sort(router.patternHandlers)
}

func (router *defaultRouter) RegisterPatternHandleFunc(priority int, pattern string, handleFunc HandleFunc, methods ...string) {
	router.RegisterPatternHandler(priority, pattern, NewHandler(handleFunc), methods...)
}

func (router *defaultRouter) Dispatch(request *core.Request) Handler {
	simpleHandlers, ok := router.simpleHandlers[request.Uri()]
	if ok {
		handler, ok := simpleHandlers[request.Method()]
		if ok {
			return handler
		}
	}

	for _, patternHandler := range router.patternHandlers {
		if patternHandler.patternMatcher.Match(request) {
			return patternHandler.handler
		}
	}

	return nil
}

func NewRouter() Router {
	return &defaultRouter{
		simpleHandlers:  make(map[string]map[string]Handler),
		patternHandlers: []*PatternHandler{},
	}
}
