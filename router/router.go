package router

import "reflect"
import "strings"
import "fmt"
import "strconv"

type Route struct {
  method reflect.Value
  args   []reflect.Type
}

var routes = make(map[string]map[string]Route)

func Register(controller interface{}) bool {
	v := reflect.ValueOf(controller)

	//非控制器或无方法则直接返回
	if v.NumMethod()==0 {
		return false
	}

	//取得控制器名称
	tmp := reflect.TypeOf(controller).String()
	module := tmp
	if strings.Contains(tmp, ".") {
		module = tmp[strings.Index(tmp, ".")+1:]

	}

	//遍历方法
    for i:= 0; i < v.NumMethod(); i++ {
    	method := v.Method(i)
    	action := v.Type().Method(i).Name

		//遍历参数
		params := make([]reflect.Type, 0, v.NumMethod())
    	for j := 0; j < method.Type().NumIn(); j++ {
      		params = append(params, method.Type().In(j))
    	}

    	if routes[module]==nil {
    		routes[module] = make(map[string]Route)
    	}

        routes[module][action] = Route{method,params}
	}

	return true
}

func Match(path string) bool {
	fields := strings.Split(path, "/")
	if len(fields) < 3 {
		return false
	}

	_, ok := routes[fields[1]][fields[2]]

	return ok
}

func Call(path string, args []string) bool {
	fields := strings.Split(path, "/")
	if len(fields) < 3 {
		return false
	}

	route, ok := routes[fields[1]][fields[2]]
	if ok==false {
		return false
	}

	argv := make([]reflect.Value, 0, len(route.args))
  	for i, t := range route.args {
    	switch t.Kind() {
    		case reflect.Int:
      			value, _ := strconv.Atoi(args[i])
      			argv = append(argv, reflect.ValueOf(value))

    		case reflect.String:
      			argv = append(argv, reflect.ValueOf(args[i]))

    		default:
      			fmt.Errorf("invalid arg type:%s", t.Kind())
      			return false
    	}
  	}

    route.method.Call(argv)

    return true
}


