package example

import (
	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"log"
	"os"
)

func Json(request *http.Request, response *http.Response) error {
	filename := "app/resources/status.json"
	pwd, _ := os.Getwd()
	log.Printf("current path is %s", pwd)
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	response.StatusOK()
	response.Json(string(content), "utf-8")
	response.End()

	return nil
}
