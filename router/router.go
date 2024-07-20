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
	ref reflect.Type
	method string
	args   []reflect.Type
	names  []string
	init bool
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

	//反射控制器类型
	refT := reflect.TypeOf(controller)

	//反射控制器值
	refV := reflect.New(refT)

	//取得控制器完整名称
	module := refT.String()

	//禁止注册控制器指针
	if strings.Contains(module, "*") {
		log.Fatal("\033[7;31;40m 控制器 ",module," 禁止注册&指针，请去掉&符号 \033[0m")
	}

	//非控制器或无方法则直接返回
	if refV.NumMethod() == 0 {
		log.Fatal("\033[7;31;40m 控制器 ",module," 无可执行结构体方法 \033[0m")
	}

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

	//打印日志
	if os.Getenv("debug")=="yes" {
		log.Println("注册控制器："+module+"{}")
	}

	maps := scan(file, module)

	//是否需要初始化
	init := refV.MethodByName("Init").IsValid()

	//遍历控制器方法
	for i := 0; i < refV.NumMethod(); i++ {
		method := refV.Method(i)
		action := refV.Type().Method(i).Name

		//忽略初始化方法
		if action=="Init" {
			continue
		}

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
			log.Fatal("\033[7;31;40m 控制器 ",module," 方法 ", action, " 参数匹配失败，请确保每个参数逗号分隔 \033[0m")
		}

		routeAction := strings.ToLower(action)
		routes.Store("/"+routeModule+"/"+routeAction, Route{refT, action, params, names, init})

		//打印日志
		if os.Getenv("debug")=="yes" {
			log.Println("\t注册方法："+action+"("+strings.Join(names, ",")+")")
		}

	}

	return true
}

//读取控制器相关源码
func read(path string) string {
	list, err := filepath.Glob(path[0:len(path)-2]+"*")
	if err != nil {
		log.Fatal("\033[7;31;40m ",path," 文件名匹配失败： ", err, " \033[0m")
	}

	var code string = ""
	for _,f := range list {
		file, err := os.Open(f)
		if err != nil {
			log.Fatal("\033[7;31;40m ",file," 文件打开失败： ", err, " \033[0m")
		}

		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("\033[7;31;40m ",file," 文件读取失败： ", err, " \033[0m")
		}

		code = code+string(content)
	}

	return code
}
//扫描控制器方法
func scan(path string, module string) map[string][]string {
	content := read(path)

	//匹配控制器方法和参数
	reg := regexp.MustCompile(`func\s*\((.+` + module + `)\s*\)\s*([A-Z_][A-Za-z0-9_]+)\s*\((.*)\)`)
	if reg == nil {
		log.Fatal("\033[7;31;40m 正则表达式编译失败 \033[0m")
	}

	maps   := map[string][]string{}
	result := reg.FindAllStringSubmatch(content, -1)
	for _, match := range result {
		receiver := match[1]
		action   := match[2]
		args     := strings.TrimSpace(match[3])

		if strings.Contains(receiver, "*")==false {
			log.Fatal("\033[7;31;40m 控制器 ",module," 方法 ", action, " 必须使用指针 *"+module+" \033[0m")
		}

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
	_, ok := routes.Load(path)
	return ok
}

//匹配路由并调用控制器方法
func Invoke(path string, args map[string]string) interface{} {
	value, ok := routes.Load(path)
	if ok == false {
		exception.New("路由地址错误:"+path, 100)
	}
	route := value.(Route)

	//检查参数并强制转换参数类型
	argvs := make([]reflect.Value, 0, len(route.args))
	for i := 0; i < len(route.names); i++ {

		name  := route.names[i]
		empty := string(name[0])=="_"
		param, ok := args[name]
		if ok == false && empty==false {
			exception.New("参数缺失:"+name, 100)
		} else if ok==false && empty==true {
			if route.args[i].Kind()==reflect.String {
				param = ""
			} else {
				param = "0"
			}
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
							return nil
		}
	}

	//生成新的零值指针
	p := reflect.New(route.ref)

	//如果需要初始化
	if route.init {
		p.MethodByName("Init").Call(make([]reflect.Value, 0, 0))
	}

	//调用目标路由方法
	result := p.MethodByName(route.method).Call(argvs)
	if len(result) == 0 {
		return nil
	}

	return result[0].Interface()
}
