package main

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/http"
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

	server := http.NewServer("0.0.0.0:4221")
	server.AddHandler("/", func(request *http.Request, response *http.Response) error {
		response.StatusOK()
		response.End()
		return nil
	})
	server.Run()
}
