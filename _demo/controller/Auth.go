package controller

import "log"
import "github.com/boyxp/nova/exception"

type Auth struct {}
func (A *Auth) Init() {
	log.Println("init ok")
	exception.New("未登录", 1001)
}
