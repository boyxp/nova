package main
//package validator

/*

init (
	validator.Register("main", func(a string, p interface{}) error {})
)

validator {
	rules map[attr]func
	
	Validate(instance interface{}, params map[string]interface{}) (result map[string]string, ok bool)
	{
		//先取规则集合
		//遍历传入参数数组
		{
			按属性找到对应函数
			传入约束和参数值
			如果报错，写入报错结果集合
		}

		如果结果集合大于0 返回数组，返回ok
	}

	getrule(instance interface{}) map[string]map[string]string {
		//先判断是否存在，没存在则解析缓存
		//返回
	}
}

*/


import "log"
import "sync"
import "errors"
import "reflect"

func main() {
	Validate(User{}, map[string]interface{}{"Mail":"abc@ccc.cc"})
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

func Validate(instance interface{}, params map[string]interface{}) (result map[string]string, ok bool) {
	result = map[string]string{}
	ok     = false

	_rules := scan(instance)

	for field, set := range _rules {
		log.Println("check:", field, set)
	}

	return result, ok
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
				rules[name][attr] = attr
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
