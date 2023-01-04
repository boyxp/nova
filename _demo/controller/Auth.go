package controller

import "log"

type Auth struct {}
func (A *Auth) Init() bool {
	log.Println("init ok")
	return true
}
