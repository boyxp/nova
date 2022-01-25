package main

import _ "nova/controller"
import "nova/router"

func main() {
	match := router.Match("/User/Login")
	if match ==true {
		router.Call("/User/Login", []string{"lee","18"})
	}
}
