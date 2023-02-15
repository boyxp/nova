package router

import "os"
import "log"
import "sync"
import "regexp"
import "strconv"
import "runtime"
import "reflect"
import "strings"
import "io/ioutil"
import "path/filepath"
import "github.com/boyxp/nova/exception"

type Route struct {
	method reflect.Value
	args   []reflect.Type
	names  []string
}

//控制器目录名称（可修改）
var controllerPathName string = "controller"

//路由规则
var routes sync.Map

//注册控制器
func Register(controller interface{}) bool {
	//取得控制器路径
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return false
	}

	//反射控制器
	v := reflect.ValueOf(controller)

	//非控制器或无方法则直接返回
	if v.NumMethod() == 0 {
		return false
	}

	//取得控制器完整名称
	module := reflect.TypeOf(controller).String()

	//取得路由模块名称
	routeModule := strings.Replace(module, "*", "", -1)
	routeModule  = strings.Replace(routeModule, controllerPathName+".", "", -1)
	routeModule  = strings.Replace(routeModule, "main.", "", -1)
	routeModule  = strings.Replace(routeModule, ".", "/", -1)
	routeModule  = strings.ToLower(routeModule)

	//取得控制器结构体名称
	if strings.Contains(module, ".") {
		module = module[strings.Index(module, ".")+1:]
	}

	maps := scan(file, module)

	//遍历控制器方法
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		action := v.Type().Method(i).Name

		//遍历方法参数取得参数类型
		params := make([]reflect.Type, 0, method.Type().NumIn())
		for j := 0; j < method.Type().NumIn(); j++ {
			params = append(params, method.Type().In(j))
		}

		//判断是否有参数名称
		names, ok := maps[action]
		if !ok {
			names = []string{}
		}

		//判断参数一致
		if len(params) != len(names) {
			panic(module + ":" + action + "参数匹配失败")
		}

		routeAction := strings.ToLower(action)
		routes.Store("/"+routeModule+"/"+routeAction, Route{method, params, names})
	}

	return true
}

//读取控制器相关源码
func read(path string) string {
	list, err := filepath.Glob(path[0:len(path)-2]+"*")
	if err != nil {
		panic(err)
	}

	var code string = ""
	for _,f := range list {
		file, err := os.Open(f)
		if err != nil {
			panic(err)
		}

		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		code = code+string(content)
	}

	return code
}
//扫描控制器方法
func scan(path string, module string) map[string][]string {
	content := read(path)

	//匹配控制器方法和参数
	reg := regexp.MustCompile(`func\s*\(.+` + module + `\s*\)\s*([A-Z][A-Za-z0-9_]+)\s*\((.*)\)`)
	if reg == nil {
		panic("MustCompile err")
	}

	maps   := map[string][]string{}
	result := reg.FindAllStringSubmatch(content, -1)
	for _, match := range result {
		action := match[1]
		args   := strings.TrimSpace(match[2])

		if len(args) == 0 {
			maps[action] = []string{}

		} else {
			sets  := []string{}
			pairs := strings.Split(args, ",")
			for i := 0; i < len(pairs); i++ {
				pairs[i] = strings.TrimSpace(pairs[i])
				pos := strings.Index(pairs[i], " ")
				if pos > -1 {
					pairs[i] = pairs[i][0:pos]
				}

				sets = append(sets, pairs[i])
			}

			maps[action] = sets
		}
	}

	return maps
}

//检查路由是否匹配
func Match(path string) bool {
	path = strings.ToLower(path)
	if strings.Contains(path, "?") {
		path = path[0:strings.Index(path, "?")]
	}

	_, ok := routes.Load(path)

	return ok
}

//匹配路由并调用控制器方法
func Invoke(path string, args map[string]string) interface{} {
	path = strings.ToLower(path)
	if strings.Contains(path, "?") {
		path = path[0:strings.Index(path, "?")]
	}

	value, ok := routes.Load(path)
	if ok == false {
		return false
	}
	route := value.(Route)

	//检查参数并强制转换参数类型
	argvs := make([]reflect.Value, 0, len(route.args))
	for i := 0; i < len(route.names); i++ {

		name := route.names[i]
		param, ok := args[name]
		if ok == false {
			exception.New("参数缺失:"+name, 100)
		}

		switch route.args[i].Kind() {
		case reflect.String:
							argvs = append(argvs, reflect.ValueOf(param))

		case reflect.Int:
							value, _ := strconv.Atoi(param)
							argvs = append(argvs, reflect.ValueOf(value))
		case reflect.Int8:
							value, _ := strconv.ParseInt(param, 10, 8)
							argvs = append(argvs, reflect.ValueOf(int8(value)))
		case reflect.Int16:
							value, _ := strconv.ParseInt(param, 10, 16)
							argvs = append(argvs, reflect.ValueOf(int16(value)))
		case reflect.Int32:
							value, _ := strconv.ParseInt(param, 10, 32)
							argvs = append(argvs, reflect.ValueOf(int32(value)))
		case reflect.Int64:
							value, _ := strconv.ParseInt(param, 10, 64)
							argvs = append(argvs, reflect.ValueOf(value))

		case reflect.Uint:
							value, _ := strconv.ParseUint(param, 10, 32)
							argvs = append(argvs, reflect.ValueOf(uint(value)))
		case reflect.Uint8:
							value, _ := strconv.ParseUint(param, 10, 8)
							argvs = append(argvs, reflect.ValueOf(uint8(value)))
		case reflect.Uint16:
							value, _ := strconv.ParseUint(param, 10, 16)
							argvs = append(argvs, reflect.ValueOf(uint16(value)))
		case reflect.Uint32:
							value, _ := strconv.ParseUint(param, 10, 32)
							argvs = append(argvs, reflect.ValueOf(uint32(value)))
		case reflect.Uint64:
							value, _ := strconv.ParseUint(param, 10, 64)
							argvs = append(argvs, reflect.ValueOf(value))

		case reflect.Bool:
							value, _ := strconv.ParseBool(param)
							argvs = append(argvs, reflect.ValueOf(value))

		case reflect.Float32:
							value, _ := strconv.ParseFloat(param, 32)
							argvs = append(argvs, reflect.ValueOf(float32(value)))
		case reflect.Float64:
							value, _ := strconv.ParseFloat(param, 64)
							argvs = append(argvs, reflect.ValueOf(value))

		default:
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
