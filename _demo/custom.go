package main

import "github.com/boyxp/nova"
import "github.com/boyxp/nova/request"
import "github.com/boyxp/nova/response"
import _ "api/controller"

func main() {
	nova.Listen(":8080").SetRequest(&request.Form{}).SetResponse(&response.Json{}).Run()
}
