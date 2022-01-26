package main

import "nova/app"
import "nova/request"
import "nova/response"
import _ "nova/controller"

func main() {
	app.New(":8080").SetRequest(&request.Form{}).SetResponse(&response.Json{}).Run()
}
