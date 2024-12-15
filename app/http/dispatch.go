package http

type Dispatcher struct {
	handlers map[string]Handler
}

func NewDispatcher() *Dispatcher {
	handlers := make(map[string]Handler)
	return &Dispatcher{handlers: handlers}
}

func (dispatcher *Dispatcher) Register(uri string, handle func(request *Request, response *Response) error) {
	dispatcher.handlers[uri] = NewHandler(handle)
}

func (dispatcher *Dispatcher) Dispatch(r *Request) Handler {
	handler, ok := dispatcher.handlers[r.Uri()]
	if ok {
		return handler
	}

	return nil
}
