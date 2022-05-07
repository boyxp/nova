package router

import "testing"
import "github.com/boyxp/nova/router"

type Hello struct {}
func (h *Hello) Hi(name string) string {
	return "hello "+name
}

func TestRegister(t *testing.T) {
	res := router.Register(&Hello{})

	if res != true {
		t.Log("路由注册失败")
		t.FailNow()
	}

	t.Log("路由注册成功")
}

func TestMatch(t *testing.T) {
	res := router.Match("/router/Hello/Hi")

	if res != true {
		t.Log("路由匹配失败")
		t.FailNow()
	}

	fail := router.Match("/router/Hello/Ok")
	if fail != false {
		t.Log("路由匹配失败")
		t.FailNow()
	}

	t.Log("路由匹配成功")
}

func TestInvoke(t *testing.T) {
	name := "eve"
	res  := router.Invoke("/router/Hello/Hi", map[string]string{"name":name})

	if res != "hello "+name {
		t.Log("路由调用失败")
		t.FailNow()
	}

	t.Log("路由调用成功")
}
