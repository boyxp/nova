package session

import "os"
import "strconv"
import "strings"
import "io/ioutil"
import "github.com/leeqvip/gophp"
import "github.com/boyxp/nova/cookie"

func Get(name string) string {
	sess := read()
	if sess==nil {
		return ""
	}

	val, ok := sess[name]
	if !ok {
		return ""
	}

	switch val.(type) {
		case string :
					value, ok := val.(string)
					if !ok {
						return ""
					}
					return value

		case int   :
				   value, ok := val.(int)
				   if !ok {
						return ""
				   }
				   return strconv.Itoa(value)
	}

	return ""
}

func Set(name string, value string) bool {
	return true
}

func ssid() string {
	ssid := cookie.Get("PHPSESSID")
	if ssid != "" {
		return ssid
	}

	ssid = cookie.Get("GOSESSID")
	return ssid
}

func path() string {
	_path := os.Getenv("session.path")
	if _path != "" {
		return _path
	}

	return "/var/lib/php/session/"
}

func read() map[string]interface{} {
	var res map[string]interface{}

	ssid := ssid()
	path := path()

	file, err := os.Open(path+"/sess_"+ssid)
	defer file.Close()
   	if err != nil {
        return res
   	}

   content, err := ioutil.ReadAll(file)
   str := string(content)
   if len(str)==0 {
   	    return res
   }

   if strings.Contains(str, "|") {
   		str = str[strings.Index(str, "|")+1:]
   }

	data, _  := gophp.Unserialize([]byte(str))
	sess, ok := data.(map[string]interface{})
	if !ok {
		return res
	}

	return sess
}
