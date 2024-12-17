package example

import (
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"log"
	"os"
)

func Json(request *core.Request, response *core.Response) error {
	filename := "resources/status.json"
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
