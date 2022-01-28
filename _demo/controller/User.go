package controller

import "fmt"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/exception"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age uint64, check bool, balance float64, num int64) string {
	if age<18 { 
		exception.New("年龄最小18岁", 101)
	}

	if age > 100 {
		panic("年龄太大")
	}

	fmt.Println(name, age, check, balance, num)

	return name;
}

func (C *User) Logout() string {
	return "bye"
}

func (C *User) private() {
	fmt.Println("private")
}
