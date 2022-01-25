package main

import _ "nova/controller"
import "nova/router"

func main() {

	match := router.Match("/Admin/Login1")
	if match ==true {
		router.Call("/Admin/Login1", []string{"aa","12"})
	}

}
