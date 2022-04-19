package cookie

import "net/http"
import "github.com/boyxp/nova/register"

func Set(name string, value string) bool {
	c := http.Cookie{
		Name    : name,
		Value   : value,
		HttpOnly: true,
		Path    : "/",
	}

	w := register.GetResponseWriter()

	http.SetCookie(w, &c)

	return true
}

func Get(name string) string {
	r := register.GetRequest()

	c, err := r.Cookie(name)
	if err == nil {
		return c.Value
	}

	return ""
}
