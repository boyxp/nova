package main

import "github.com/boyxp/nova"
import "github.com/boyxp/nova/database"
import _ "api/controller"

func main() {
	database.Register("database", "dev_xinzhanghu", "root:123456@tcp(localhost:3306)/dev_xinzhanghu")
	nova.Listen("8080").Run()
}
