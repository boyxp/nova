package main
//package validator

import "log"
import "sync"
import "errors"
import "reflect"

func main() {
	result := Validate(User{}, map[string]interface{}{"Mail":"abc@ccc.cc"})
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
	Temp string `min:"1" max:"2" length:"10"`
	Empty string
}

func init() {
	Register("mail", func(set string, param interface{}) error{
		log.Println(param)
		return errors.New("不是邮箱格式")
	})

	Register("url", func(set string, param interface{}) error{
		return errors.New("不符合要求")
	})

	Register("ip", func(set string, param interface{}) error{
		return errors.New("不符合要求")
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
		return errors.New("不符合要求")
	})

	Register("max", func(set string, param interface{}) error{
		return errors.New("不符合要求")
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

	rules := scan(instance)

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
					err = err+_res.Error()
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
    	log.Println("cache:yes")
        return value.(map[string]map[string]string)
    }


	rules = map[string]map[string]string{}

	for i := 0; i < ref.NumField(); i++ {
		field      := ref.Field(i)
		name       := field.Name
		tag        := field.Tag
		rules[name] = map[string]string{}

		if tag=="" {
			continue
		}

		_tag := string(tag)
		attrs.Range(func(key, value interface{}) bool {
			attr := key.(string)
			if _tag== attr {
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
