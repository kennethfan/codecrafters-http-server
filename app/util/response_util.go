package util

import (
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"github.com/kennethfan/codecrafters-http-server/http/core"
)

func Response200(response *core.Response) {
	response.SetStatus(common.StatusOK)
	response.SetMessage(common.MessageOk)
}

func ResponseSuccess(response *core.Response, message string) {
	response.SetStatus(common.StatusOK)
	response.SetMessage(message)
}

func ResponseError(response *core.Response, status uint32, message string) {
	response.SetStatus(status)
	response.SetMessage(message)
}
func Response404(response *core.Response) {
	response.SetStatus(common.StatusNotFound)
	response.SetMessage(common.MessageNotFound)
}
func Response403(response *core.Response) {
	response.SetStatus(common.StatusForbidden)
	response.SetMessage(common.MessageForbidden)
}
func Response500(response *core.Response) {
	response.SetStatus(common.StatusInternalServerError)
	response.SetMessage(common.MessageInternalServerError)
}

func ResponseSuccessHtml(response *core.Response, encoding string, content []byte) {
	Response200(response)
	response.SetBody(common.ContentTypeHtml, encoding, content)
}

func ResponseSuccessPlain(response *core.Response, encoding string, content []byte) {
	Response200(response)
	response.SetBody(common.ContentTypePlain, encoding, content)
}

func ResponseSuccessJson(response *core.Response, encoding string, content []byte) {
	Response200(response)
	response.SetBody(common.ContentTypeJson, encoding, content)
}
