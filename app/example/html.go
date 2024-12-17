package example

import (
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"log"
	"os"
)

func Html(request *core.Request, response *core.Response) error {
	filename := "resources/index.html"
	pwd, _ := os.Getwd()
	log.Printf("current path is %s", pwd)
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	response.StatusOK()
	response.Html(string(content), "utf-8")
	response.End()

	return nil
}
