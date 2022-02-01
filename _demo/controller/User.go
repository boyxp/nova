package controller

import "fmt"
import "log"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/exception"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age uint64, check bool, balance float64, num int64) interface{} {
	log.Println("姓名：",name,"年龄：",age,"检查：",check,"余额：",balance,"数量：",num)

	if age<18 {
		exception.New("年龄最小18岁", 101)
	}

	if age > 100 {
		panic("年龄参数异常")
	}

	//返回字符串
	//return name

	//返回int64
	//return balance

	//返回布尔值
	//return check

	//返回切片
	//return []string{name,name,name}

	//返回map
	//return map[string]interface{}{"Name":name,"Age":age}

	//返回复杂结果集
	list := []interface{}{map[string]interface{}{"Name":name,"Age":age},map[string]interface{}{"Name":name,"Age":age}}
	return map[string]interface{}{"total":len(list),"list":list}
}

func (C *User) Logout() {
	fmt.Println("bye")
}

func (C *User) private() {
	fmt.Println("private")
}
