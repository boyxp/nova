package controller

import "os"
import "fmt"
import "log"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/exception"
import "github.com/boyxp/nova/response"

import "api/model"

func init() {
   router.Register(User{})
}

type User struct {
	Auth
}

func (C User) Hello() map[string]string {
	return map[string]string{"Item0":"0、当前项目目录为 _demo","Item1":"1、端口和数据库配置在 .env","Item2":"2、控制器目录为controller","Item3":"3、模型目录为model","Item4":"4、进程管理使用 sh manage.sh","Item5":"5、更多示例见User控制器",}
}

func (C User) Login(name string, age uint64, check bool, balance float64, num int64, portrait string) any {
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
	list := []any{map[string]any{"Name":name,"Age":age,"Portrait":portrait},map[string]any{"Name":name,"Age":age,"Portrait":portrait}}
	return map[string]any{"total":len(list),"list":list}
}

//返回其他结构
func (C User) Logout() map[string]string {
	fmt.Println("bye")
	return map[string]string{"res":"bye"}
}

//文件下载
func (C User) Download() {
	d := response.Download{}
	d.Render("文件.log")
}

//重定向
func (C User) Jump() {
	d := response.Redirect{}
	d.Render("https://www.baidu.com")
}

//数据库添加操作
func (C User) Add() map[string]any {
	user_id := model.User.Insert(map[string]string{"user":"xiaoming","password":"123"})
	return map[string]any{"user_id":user_id}
}

//数据库列表读取
func (C User) List() []map[string]string {
	list := model.User.Where("user", "xiaoming").Select()
	return list
}

//小写开头私有方法不可通过路由访问
func (C User) private() {
	fmt.Println("private")
}
