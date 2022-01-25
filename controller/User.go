package controller

import "fmt"
import "hello/router"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age int) {
	fmt.Println(name, age, "ok")
}

func (C *User) Logout() {
	fmt.Println("bye")
}
