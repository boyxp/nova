package controller

import "fmt"
import "nova/router"
import "time"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age int) string {
	time.Sleep(time.Duration(10 * time.Second))
	return name;
}

func (C *User) Logout() string {
	return "bye"
}

func (C *User) private() {
	fmt.Println("private")
}
