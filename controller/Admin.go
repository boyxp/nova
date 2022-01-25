package controller

import "fmt"
import "nova/router"

func init() {
   router.Register(&Admin{})
}

type Admin struct {}
func (C *Admin) Login1(name string, age int) {
	fmt.Println(name, age, "o----k")
}

func (C *Admin) Logout() {
	fmt.Println("bye")
}
