package controller

import "fmt"
import "nova/router"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age int) {
	fmt.Println(name, age, "hello")
}

func (C *User) Logout() {
	fmt.Println("bye")
}

func (C *User) private() {
	fmt.Println("private")
}
