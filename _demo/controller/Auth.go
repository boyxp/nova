package controller

import "log"

type Auth struct {}

func (A *Auth) Init() {
	//可以在这里检查登录状态
	log.Println("init ok")
}
