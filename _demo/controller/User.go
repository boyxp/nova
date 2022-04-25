package controller

import "os"
import "fmt"
import "log"
import "net/http"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/exception"
import "github.com/boyxp/nova/response"

import "api/model"

func init() {
   router.Register(&User{})
}

type User struct {}
func (C *User) Login(name string, age uint64, check bool, balance float64, num int64, portrait string) interface{} {
	log.Println("姓名：",name,"年龄：",age,"检查：",check,"余额：",balance,"数量：",num)

	if age<18 {
		//逻辑异常
		exception.New("年龄最小18岁", 1001)
	}

	if age > 100 {
		//非逻辑异常
		panic("年龄参数异常")
	}

	//文件大小检查
	file, err := os.Stat(portrait)
	if err == nil {
		if file.Size() > 1024000 {
			exception.New("头像不得大于1m", 1002)
		}
	}



	//文件格式检查
	fp, err := os.Open(portrait)
	if err != nil {
			panic(err)
	}
	defer fp.Close()

	buffer := make([]byte, 512)
	_, err1 := fp.Read(buffer)
	if err1 != nil {
		panic(err1)
	}

	contentType := http.DetectContentType(buffer)

	if contentType!="image/jpeg" {
		exception.New("头像不是图片格式", 1003)
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
	list := []interface{}{map[string]interface{}{"Name":name,"Age":age,"Portrait":portrait},map[string]interface{}{"Name":name,"Age":age,"Portrait":portrait}}
	return map[string]interface{}{"total":len(list),"list":list}
}

func (C *User) Logout() map[string]string {
	fmt.Println("bye")
	return map[string]string{"res":"bye"}
}

func (C *User) Download() {
	d := response.Download{}
	d.Render("文件.log")
}

func (C *User) Jump() {
	d := response.Redirect{}
	d.Render("https://www.baidu.com")
}

func (C *User) Add() map[string]interface{} {
	user_id := model.User().Insert(map[string]interface{}{"user_name":"xiaoming","password":"123"})
	return map[string]interface{}{"user_id":user_id}
}

func (C *User) List() []map[string]string {
	list := model.User().Select()
	return list
}

func (C *User) private() {
	fmt.Println("private")
}
