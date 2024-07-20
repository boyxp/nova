package router

import "testing"

type Auth struct {
	Uid string
}

func (A *Auth) Init() {
	A.Uid = "1024"
}

type Hello struct {
	Auth
}

func (h *Hello) Hi(name string) string {
	return "hello "+name+h.Uid
}



func TestRegister(t *testing.T) {
	res := Register(Hello{})

	if res != true {
		t.Log("路由注册失败")
		t.FailNow()
	}

	t.Log("路由注册成功")
}

func TestMatch(t *testing.T) {
	res := Match("/router/Hello/Hi")

	if res != true {
		t.Log("路由匹配失败")
		t.FailNow()
	}

	fail := Match("/router/Hello/Ok")
	if fail != false {
		t.Log("路由匹配失败")
		t.FailNow()
	}

	t.Log("路由匹配成功")
}

func TestInvoke(t *testing.T) {
	name := "eve"
	res  := Invoke("/router/Hello/Hi", map[string]string{"name":name})

	if res != "hello "+name+"1024" {
		t.Log("路由调用失败")
		t.FailNow()
	}

	t.Log("路由调用成功")
}
