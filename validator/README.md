## 验证器

### 用法示例

```go
package main

import "log"
import "github.com/boyxp/nova/validator"

func main() {
        result := validator.Validate(User{}, map[string]interface{}{
                "UMail":"a.b-c+d@efg.",
                "Num":9,
                "UUrl":"htt://www.com:8080/a/b/c?a=b&c=d#1111",
                "IP":"a.0.0.0",
                "Phone":"01-12345678999",
                "Mobi":"188888888880",
                "Card":"12345678",
                "Age":160,
                "Height":175,
                "Color":"black",
        })

        for field, err := range result {
                log.Println("参数：",field, "错误：",err)
        }
}

type User struct {
        UMail string `mail`
        UUrl string `url`
        IP string `ip`
        Mobi string `mobile`
        Phone string `tel`
        Num string `min:"10" max:"20" length:"10"`
        Addr string
        Card string `len:"7"`
        Age int `gt:"18" lt:"50" eq:"40"`
        Height int `gte:"160" lte:"170" ne:"175"`
        Color string `in:"red,green,blue"`
}
```

### 注册自定义规则

set为设定值，param为实际参数

```go
        validator.Register("must", func(set string, param interface{}) error{
                _param, ok := param.(string)
                if !ok {
                        return errors.New("参数类型错误")
                }

                if(set==_param) {
                        return nil
                }

                return errors.New("值必须为"+set)
        })

        type Book struct {
                Sort string `must:"novel"`
        }

        validator.Validate(Book{}, map[string]interface{}{
                "Sort":"story",
        })
```

