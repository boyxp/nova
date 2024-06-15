package main

import "log"
import "net/http"
import "github.com/boyxp/nova"
import _ "api/controller"

func main() {
	nova.Listen().Use(logger).Run()
}

func logger(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      log.Println("logger:", r.URL.Path)
      next.ServeHTTP(w, r)
   })
}
