package main

import "os"
import "github.com/boyxp/nova"
import "github.com/boyxp/nova/database"
import _ "api/controller"

func main() {
	database.Register("database", os.Getenv("database.dbname"), os.Getenv("database.dsn"))
	nova.Listen(os.Getenv("port")).Run()
}
