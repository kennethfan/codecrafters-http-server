package main

import (
	"fmt"
	"github.com/kennethfan/codecrafters-http-server/example"
	"github.com/kennethfan/codecrafters-http-server/http"
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"github.com/kennethfan/codecrafters-http-server/http/component"
	"github.com/kennethfan/codecrafters-http-server/http/component/builtin"
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"github.com/kennethfan/codecrafters-http-server/util"
	"log"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//

	server := http.NewServer()
	//server := http.NewServer(http.WithAddress("0.0.0.0:4221"), http.WithRouter(http.NewRouter()))
	router := server.Router()
	router.RegisterSimpleHandleFunc("/", func(request *core.Request, response *core.Response) error {
		util.Response200(response)
		return nil
	})
	router.RegisterSimpleHandleFunc("/index.html", example.Html, common.HttpMethodGet)
	router.RegisterSimpleHandleFunc("/status.json", example.Json, common.HttpMethodGet)
	mimeTypes, err := component.MimeTypeFromFile("resources/mime.types")
	if err != nil {
		log.Println("build mimetypes error, file=resources/mime.types")
	} else {
		router.RegisterPatternHandler(1, "/static/.+", builtin.NewStaticHandler("/static/", "static/", mimeTypes), common.HttpMethodGet)
	}

	middlewares := server.Middlewares()
	middlewares.RegisterPatternMiddleware(1, "/.*", &builtin.LogMiddleware{})
	server.Start()
}
