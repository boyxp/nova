package router

import "testing"
import "log"

type Auth struct {
	Uid string
}

func (A *Auth) Init() {
	A.Uid = "1024"
}

type Hello struct {
	Auth
}

func (h *Hello) Hi(name string, _remark string, _num int) string {
	log.Println(_remark, _num)
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
	res := Match("/router/hello/hi")

	if res != true {
		t.Log("路由匹配失败")
		t.FailNow()
	}

	fail := Match("/router/hello/ok")
	if fail != false {
		t.Log("路由匹配失败")
		t.FailNow()
	}

	t.Log("路由匹配成功")
}

func TestInvoke(t *testing.T) {
	name := "eve"
	res  := Invoke("/router/hello/hi", map[string]string{"name":name, "_remark":"ok", "_num":"1024"})

	if res != "hello "+name+"1024" {
		t.Log("路由调用失败")
		t.FailNow()
	}

	t.Log("路由调用成功")
}
