package main

import "os"
import "log"
import "net/http"
import "github.com/boyxp/nova"
import "github.com/boyxp/nova/database"
import _ "api/controller"

func main() {
	database.Register("database", os.Getenv("database.dbname"), os.Getenv("database.dsn"))
	nova.Listen(os.Getenv("port")).Use(logger).Run()
}

func logger(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      log.Println("logger:", r.URL.Path)
      next.ServeHTTP(w, r)
   })
}
