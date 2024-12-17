package example

import (
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"os"
)

func Html(request *core.Request, response *core.Response) error {
	filename := "resources/index.html"
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	response.SetStatus(common.StatusOK)
	response.SetMessage(common.MessageOk)
	response.SetBody("text/html", "utf-8", content)

	return nil
}
