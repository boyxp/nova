package main

import "log"
import "os"
import "github.com/boyxp/nova/database"
import _ "github.com/joho/godotenv/autoload"
import "api/model"

func main() {
	database.Register("database", os.Getenv("database.dbname"), os.Getenv("database.dsn"))

	list := model.User().Select()
	log.Println(list)
}
