package controller

func (h Hello) Haha(name string) map[string]string {
	return map[string]string{"name":"haha "+name}
}
