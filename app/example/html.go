package example

import (
	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"log"
	"os"
)

func Html(request *http.Request, response *http.Response) error {
	filename := "app/resources/index.html"
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
