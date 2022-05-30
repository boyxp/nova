package session

import "os"
import "sync"
import "strconv"
import "strings"
import "io/ioutil"
import "github.com/techoner/gophp"
import "github.com/boyxp/nova/cookie"

var cache sync.Map
var path string = "/var/lib/php/session/"

func Config(set string) bool {
	path = set
	return true
}

func Get(name string) string {
	ssid := cookie.Get("PHPSESSID")
	if ssid == "" {
		return ""
	}

	file, err := os.Open(path+"/sess_"+ssid)
	defer file.Close()
   	if err != nil {
        return ""
   	}

   content, err := ioutil.ReadAll(file)
   str := string(content)
   if len(str)==0 {
   	    return ""
   }

   if strings.Contains(str, "|") {
   		str = str[strings.Index(str, "|")+1:]
   }

	data, _  := gophp.Unserialize([]byte(str))
	sess, ok := data.(map[string]interface{})
	if !ok {
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
