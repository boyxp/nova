package main

import "fmt"
import "github.com/boyxp/nova/time"

func main() {
        fmt.Println(time.Date("Y-m-d H:i:s"))
        fmt.Println(time.Strtotime("-30 days"))
        fmt.Println(time.Date("Y-m-d", time.Strtotime("-3 months")))

        fmt.Println(time.Now().Format("2006_01_02-15_04_05"))
}
