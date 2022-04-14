package main

import "github.com/boyxp/nova"
import _ "api/controller"

func main() {
	nova.Listen("8080").Run()
}
