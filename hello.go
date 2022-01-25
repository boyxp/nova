package main

import _ "hello/controller"
import "hello/router"

func main() {

	match := router.Match("/Admin/Login1")
	if match ==true {
		router.Call("/Admin/Login1", []string{"aa","12"})
	}

}
