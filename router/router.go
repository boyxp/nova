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

//注册控制器
func Register(controller interface{}) bool {
	v := reflect.ValueOf(controller)

	//非控制器或无方法则直接返回
	if v.NumMethod()==0 {
		return false
	}

	//取得控制器名称
	module := reflect.TypeOf(controller).String()
	if strings.Contains(module, ".") {
		module = module[strings.Index(module, ".")+1:]
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

//检查路由是否匹配
func Match(path string) bool {
	if strings.Contains(path, "?") {
		path = path[0:strings.Index(path, "?")]
	}

	fields := strings.Split(path, "/")
	if len(fields) < 3 {
		return false
	}

	_, ok := routes[fields[1]][fields[2]]

	return ok
}

//调用路由方法
func Call(path string, args []string) bool {
	fields := strings.Split(path, "/")
	if len(fields) < 3 {
		return false
	}

	route, ok := routes[fields[1]][fields[2]]
	if ok==false {
		return false
	}

	argvs := make([]reflect.Value, 0, len(route.args))
  for i, t := range route.args {
    	switch t.Kind() {
    		case reflect.Int:
      									value, _ := strconv.Atoi(args[i])
      									argvs     = append(argvs, reflect.ValueOf(value))

    		case reflect.String:
      									argvs = append(argvs, reflect.ValueOf(args[i]))

    		default:
      			fmt.Errorf("invalid arg type:%s", t.Kind())
      			return false
    		}
  	}

    route.method.Call(argvs)

    return true
}
