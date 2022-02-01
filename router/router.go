package router

import "reflect"
import "strings"
import "log"
import "strconv"

type Route struct {
  method reflect.Value
  args   []reflect.Type
}

//路由规则
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

	//遍历控制器方法
  for i:= 0; i < v.NumMethod(); i++ {
    	method := v.Method(i)
    	action := v.Type().Method(i).Name

			//遍历方法参数取得参数类型
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

//匹配路由并调用控制器方法
func Invoke(path string, args []string) interface{} {
	if strings.Contains(path, "?") {
		path = path[0:strings.Index(path, "?")]
	}

	fields := strings.Split(path, "/")
	if len(fields) < 3 {
		return false
	}

	route, ok := routes[fields[1]][fields[2]]
	if ok==false {
		return false
	}

	//判断POST参数个数是否少于方法参数
	if len(route.args)>len(args) {
		return false
	}

	//强制转换参数类型
	argvs := make([]reflect.Value, 0, len(route.args))
  for i:=0;i<len(route.args);i++ {

    	switch route.args[i].Kind() {
    		case reflect.String :
      											 argvs = append(argvs, reflect.ValueOf(args[i]))

      	case reflect.Int    :
      											 value, _ := strconv.Atoi(args[i])
        										 argvs     = append(argvs, reflect.ValueOf(value))
    		case reflect.Int8   :
      											 value, _ := strconv.ParseInt(args[i], 10, 8)
      											 argvs     = append(argvs, reflect.ValueOf(int8(value)))
    		case reflect.Int16  :
      											 value, _ := strconv.ParseInt(args[i], 10, 16)
      											 argvs     = append(argvs, reflect.ValueOf(int16(value)))
    		case reflect.Int32  :
      											 value, _ := strconv.ParseInt(args[i], 10, 32)
      											 argvs     = append(argvs, reflect.ValueOf(int32(value)))
    		case reflect.Int64  :
      											 value, _ := strconv.ParseInt(args[i], 10, 64)
      											 argvs     = append(argvs, reflect.ValueOf(value))

      	case reflect.Uint   :
      											 value, _ := strconv.ParseUint(args[i], 10, 32)
      											 argvs     = append(argvs, reflect.ValueOf(uint(value)))
    		case reflect.Uint8  :
      											 value, _ := strconv.ParseUint(args[i], 10, 8)
      											 argvs     = append(argvs, reflect.ValueOf(uint8(value)))
      	case reflect.Uint16 :
      											 value, _ := strconv.ParseUint(args[i], 10, 16)
      											 argvs     = append(argvs, reflect.ValueOf(uint16(value)))
      	case reflect.Uint32 :
      											 value, _ := strconv.ParseUint(args[i], 10, 32)
      											 argvs     = append(argvs, reflect.ValueOf(uint32(value)))
    		case reflect.Uint64 :
      											 value, _ := strconv.ParseUint(args[i], 10, 64)
      											 argvs     = append(argvs, reflect.ValueOf(value))

    		case reflect.Bool   :
      											 value, _ := strconv.ParseBool(args[i])
      											 argvs     = append(argvs, reflect.ValueOf(value))

				case reflect.Float32:
        										 value, _ := strconv.ParseFloat(args[i], 32)
        										 argvs     = append(argvs, reflect.ValueOf(float32(value)))
    		case reflect.Float64:
      											 value, _ := strconv.ParseFloat(args[i], 64)
      											 argvs     = append(argvs, reflect.ValueOf(value))

    		default							:
    												 log.Printf("Unsupported argument type:%s", route.args[i].Kind())
      											 return false
    		}
  	}

    result := route.method.Call(argvs)
    if len(result) == 0 {
    	return nil
    }

    return result[0].Interface()
}
