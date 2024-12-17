package builtin

import (
	"github.com/kennethfan/codecrafters-http-server/http/component"
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"log"
)

type LogMiddleware struct {
}

func (logMiddleware *LogMiddleware) DoChain(request *core.Request, response *core.Response, chain component.MiddlewareChain) (bool, error) {
	ok, err := chain.DoChain(request, response)
	log.Printf("%s completed result=%v, request=%v, response=%v, error=%v", request.Uri(), ok, request, response, err)
	return ok, err
}
