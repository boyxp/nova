package controller

import "fmt"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/exception"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age uint64, check bool, balance float64, num int64) interface{} {
	fmt.Println("姓名：",name,"年龄：",age,"检查：",check,"余额：",balance,"数量：",num)

	if age<18 {
		exception.New("年龄最小18岁", 101)
	}

	if age > 100 {
		panic("年龄参数异常")
	}

	return map[string]interface{}{"Name":name}
}

func (C *User) Logout() {
	fmt.Println("bye")
}

func (C *User) private() {
	fmt.Println("private")
}
