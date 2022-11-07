package main
//package validator

import "log"
import "net"
import "sync"
import "regexp"
import "errors"
import "reflect"
import "strconv"

func main() {
	result := Validate(User{}, map[string]interface{}{"Mail":"a.b-c+d@efg.cn", "Temp":11,"Url":"https://www.com:8080/a/b/c?a=b&c=d#1111","Ip":"0.0.0.0"})
	for f,e := range result {
		log.Println("参数：",f,"错误：",e)
	}
}

type User struct {
	Mail string `mail`
	Url string `url`
	Ip string `ip`
	Date string `date`
	Mobile string `mobile`
	Tel string `tel`
	Temp string `min:"10" max:"20" length:"10"`
	Empty string
}

func init() {
	Register("mail", func(set string, param interface{}) error{
		_param, ok := param.(string)
		if !ok {
			return errors.New("参数类型错误")
		}

		pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
		reg     := regexp.MustCompile(pattern)
		res     := reg.MatchString(_param)

		if(res==true) {
			return nil
		}

		return errors.New("不是邮箱格式")
	})

	Register("url", func(set string, param interface{}) error{
		_param, ok := param.(string)
		if !ok {
			return errors.New("参数类型错误")
		}

		reg    := regexp.MustCompile("(http|https):\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?")
		result := reg.FindAllStringSubmatch(_param, -1)
		if result == nil {
			return errors.New("不是url格式")
		}

		return nil
	})

	Register("ip", func(set string, param interface{}) error{
		_param, ok := param.(string)
		if !ok {
			return errors.New("参数类型错误")
		}

		ip := net.ParseIP(_param)

		if ip != nil && ip.To4() != nil {
			return nil
		}

		return errors.New("不是ip格式")
	})

	Register("tel", func(set string, param interface{}) error{
		return errors.New("不符合要求")
	})

	Register("mobile", func(set string, param interface{}) error{
		return errors.New("不符合要求")
	})

	Register("domain", func(set string, param interface{}) error{
		return errors.New("不符合要求")
	})

	Register("min", func(set string, param interface{}) error{
		_set, err := strconv.Atoi(set)
		if err != nil {
			return errors.New("min设置类型错误")
		}

		_param, ok := param.(int)
		if !ok {
			return errors.New("参数类型错误")
		}

		if _param >= _set {
			return nil
		}

		return errors.New("不可小于"+set)
	})

	Register("max", func(set string, param interface{}) error{
		_set, err := strconv.Atoi(set)
		if err != nil {
			return errors.New("max设置类型错误")
		}

		_param, ok := param.(int)
		if !ok {
			return errors.New("参数类型错误")
		}

		if _param <= _set {
			return nil
		}

		return errors.New("不可大于"+set)
	})

	//"date":"date",len":"length","length":"length","gt":"gt","gte":"gte","lt":"lt","lte":"lte","ne":"ne","size":"size"
}

//属性
var attrs sync.Map

//验证规则缓存
var cache sync.Map

func Register(attr string, call func(set string, param interface{}) error) bool {
	attrs.Store(attr, call)
	return true
}

func Validate(instance interface{}, params map[string]interface{}) map[string]string {
	result := map[string]string{}
	rules  := scan(instance)

	if len(rules)==0 {
		return result
	}

	for field, param := range params {
		err := ""
		if sets, ok := rules[field];ok {
			for attr,set := range sets {
				_call, _ := attrs.Load(attr)
				_func    := _call.(func(set string, param interface{}) error)
				_res     := _func(set, param)

				if _res!=nil {
					err = err+_res.Error()+";"
				}
			}
		}

		if len(err)>0 {
			result[field] = err
		}
	}

	return result
}

func scan(instance interface{}) (rules map[string]map[string]string) {
	ref := reflect.TypeOf(instance)
    key := ref.String()

	value, ok := cache.Load(key)
    if ok {
        return value.(map[string]map[string]string)
    }


	rules = map[string]map[string]string{}

	for i := 0; i < ref.NumField(); i++ {
		field      := ref.Field(i)
		name       := field.Name
		tag        := field.Tag

		if tag=="" {
			continue
		}

		rules[name] = map[string]string{}

		_tag := string(tag)
		attrs.Range(func(key, value interface{}) bool {
			attr := key.(string)
			if _tag==attr {
				rules[name][attr] = ""
				return false
			}

			if set, ok := tag.Lookup(attr);ok {
				rules[name][attr] = set
			}

			return true
		})
	}

	cache.Store(key, rules)

	return rules
}
