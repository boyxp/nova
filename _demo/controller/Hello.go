package controller

import "github.com/boyxp/nova/router"
func init() {
   router.Register(&Hello{})
}

type Hello struct {
	Auth
}
func (h *Hello) Hi(name string) map[string]string {
	return map[string]string{"name":"hello "+name}
}
