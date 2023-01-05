package cookie

import "os"
import "strconv"
import "net/http"
import "github.com/boyxp/nova/register"

func Set(name string, value string) bool {
	HttpOnly := true
	if read("HttpOnly")=="false" {
		HttpOnly = false
	}

	Secure := true
	if read("Secure")=="false" {
		Secure = false
	}

	Path := "/"
	if read("Path")!="" {
		Path = read("Path")
	}

	Domain := read("Domain")

	MaxAge := 86400
	if read("MaxAge")!="" {
		v, err := strconv.Atoi(read("MaxAge"))
		if err==nil {
			MaxAge = v
		}
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

func read(key string) string {
	return os.Getenv("cookie."+key)
}
