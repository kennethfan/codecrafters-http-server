package http

type Handler interface {
	Handle(*Request, *Response) error
}

type HandlerImpl struct {
	f func(request *Request, response *Response) error
}

func (handler *HandlerImpl) Handle(request *Request, response *Response) error {
	err := handler.f(request, response)
	return err
}

func NewHandler(handle func(r *Request, response *Response) error) Handler {
	return &HandlerImpl{f: handle}
}
