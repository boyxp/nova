package controller

import "log"

type Auth struct {}
func (A *Auth) Init() {
	log.Println("init ok")
}
