package main

import "os"
import "github.com/boyxp/nova"
import _ "api/controller"

func main() {
	nova.Listen(os.Getenv("port")).Run()
}
