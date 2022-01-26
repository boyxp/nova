package controller

import "fmt"
import "github.com/boyxp/nova/router"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age int, check bool, balance float64) string {
	fmt.Println(name, age, check, balance)
	return name;
}

func (C *User) Logout() string {
	return "bye"
}

func (C *User) private() {
	fmt.Println("private")
}
