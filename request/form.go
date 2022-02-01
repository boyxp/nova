package request

import "net/http"
import "io/ioutil"
import "strings"

type Form struct {}

func (F *Form) Parse(r *http.Request)map[string]string {
	params := map[string]string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return params
	}

	pairs := strings.Split(string(body), "&")

	for i:=0;i<len(pairs);i++ {
		pos := strings.Index(pairs[i], "=")
		key := pairs[i][0:pos]
		val := pairs[i][pos+1:]
		params[key] = val
	}

	return params
}
