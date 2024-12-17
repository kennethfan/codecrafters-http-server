package example

import (
	"github.com/kennethfan/codecrafters-http-server/http/common"
	"github.com/kennethfan/codecrafters-http-server/http/core"
	"os"
)

func Json(request *core.Request, response *core.Response) error {
	filename := "resources/status.json"
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	response.SetStatus(common.StatusOK)
	response.SetMessage(common.MessageOk)
	response.SetBody("application/json", "utf-8", content)

	return nil
}
