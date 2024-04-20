package main

import "log"
import _ "github.com/joho/godotenv/autoload"
import "api/model"

func main() {
	list := model.User.Select()
	log.Println(list)
}
