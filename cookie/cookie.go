package cookie

import "net/http"
import "github.com/boyxp/nova/register"

var config map[string]interface{}
func Config(set map[string]interface{}) bool {
	config = set
	return true
}

func Set(name string, value string) bool {
	HttpOnly,ok := read("HttpOnly").(bool)
	if !ok {
		HttpOnly = true
	}

	Secure,ok := read("Secure").(bool)
	if !ok {
		Secure = false
	}

	Path,ok := read("Path").(string)
	if !ok {
		Path = "/"
	}

	Domain,ok := read("Domain").(string)
	if !ok {
		Domain = ""
	}

	MaxAge,ok := read("MaxAge").(int)
	if !ok {
		MaxAge = 1200
	}

	c := http.Cookie{
		Name    : name,
		Value   : value,
		HttpOnly: HttpOnly,
		Secure  : Secure,
		Path    : Path,
		Domain  : Domain,
		MaxAge  : MaxAge,
	}

	w := register.GetResponseWriter()
	if w==nil {
		return false
	}

	http.SetCookie(w, &c)

	return true
}

func Get(name string) string {
	r := register.GetRequest()
	if r==nil {
		return ""
	}

	c, err := r.Cookie(name)
	if err == nil {
		return c.Value
	}

	return ""
}

func read(key string) interface{} {
	value, ok := config[key]

	if ok {
		return value
	}

	return nil
}
