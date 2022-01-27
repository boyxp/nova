package nova

import "log"
import "os"
import "net/http"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/request"
import "github.com/boyxp/nova/response"
import "github.com/fvbock/endless"

func New(port string) *App {
	return &App{port, &request.Form{}, &response.Json{}}
}

type App struct {
	Port string
	Request request.Interface
	Response response.Interface
}

func (A *App) Run() {
	err := endless.ListenAndServe(A.Port, A)
	if err != nil {
		log.Println(err)
	}
	log.Println("Server stopped")

	os.Exit(0)
}

func (A *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	match := router.Match(r.RequestURI)
	if match ==true {
		res := router.Call(r.RequestURI, []string{"lee","1","aa","10.11","1"})
		log.Println(res)
	}
}

func (A *App) SetRequest(req request.Interface) *App {
	A.Request = req
	return A
}

func (A *App) SetResponse(res response.Interface) *App {
	A.Response = res
	return A
}
