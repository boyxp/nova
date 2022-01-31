package request

import "net/http"
import "io/ioutil"
import "strings"

type Form struct {}

func (F *Form) Parse(r *http.Request)[]string {
	params := []string{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return params
	}

	pairs := strings.Split(string(body), "&")

	for i:=0;i<len(pairs);i++ {
		pos := strings.Index(pairs[i], "=")+1
		val := pairs[i][pos:]

		params = append(params, val)
	}

	return params
}
