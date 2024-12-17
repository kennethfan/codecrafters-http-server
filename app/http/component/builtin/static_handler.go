package builtin

import (
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"github.com/kennethfan/codecrafters-http-server/http/component"
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"github.com/kennethfan/codecrafters-http-server/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type StaticHandler struct {
	prefix    string
	rootPath  string
	mimeTypes component.MimeTypes
}

func (handler StaticHandler) Handle(request *core.Request, response *core.Response) error {
	if len(request.Uri()) <= len(handler.prefix) {
		util.Response404(response)
		return nil
	}

	if !strings.EqualFold(request.Uri()[0:len(handler.prefix)], handler.prefix) {
		util.Response404(response)
		return nil
	}
	path := filepath.Join(handler.rootPath, request.Uri()[len(handler.prefix):])
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			util.Response404(response)
			return nil
		}
		if os.IsPermission(err) {
			util.Response403(response)
		}

		log.Printf("stat file %s error: %v", path, err)
		util.Response500(response)
		return nil
	}

	if info.IsDir() {
		util.Response403(response)
		return nil
	}
	body, err := os.ReadFile(path)
	if err != nil {
		log.Printf("read file %s error: %v", path, err)
		util.Response500(response)
		return nil
	}

	util.Response200(response)
	contentType := common.ContentTypeOctetStream
	ext := filepath.Ext(path)
	if ext != "" {
		contentType = handler.mimeTypes.ContentType(ext[1:])
	}
	if contentType != common.ContentTypeOctetStream {
		response.SetHeader(common.HeaderContentType, contentType)
		response.SetBody(contentType, "", body)
		return nil
	}

	response.SetHeader(common.HeaderContentDisposition, "attachment; filename="+info.Name())
	response.SetBody(contentType, "", body)
	return nil
}

func NewStaticHandler(prefix string, rootPath string, mimeTypes component.MimeTypes) *StaticHandler {
	return &StaticHandler{
		prefix:    prefix,
		rootPath:  rootPath,
		mimeTypes: mimeTypes,
	}
}
