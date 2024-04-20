## 常用时间

### 使用示例

```
package main

import "fmt"
import "github.com/boyxp/nova/time"

func main() {
        //直接格式化当前时间
        fmt.Println(time.Date("Y-m-d H:i:s"))

        //口语时间转换为标准时间
        fmt.Println(time.Strtotime("-30 days"))
        fmt.Println(time.Strtotime("next month"))

        //转换和格式化组合使用
        fmt.Println(time.Date("Y-m-d", time.Strtotime("-3 months +10 days")))

        //以Now()为入口使用系统time包的其他原生方法，避免引用包混淆
        fmt.Println(time.Now().Format("2006_01_02-15_04_05"))
}
```
