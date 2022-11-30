package session

import "log"

import "os"
import "net"
import "runtime"
import "strconv"
import "strings"
import "net/http"
import "io/ioutil"
import "crypto/md5"
import "encoding/hex"
import "github.com/leeqvip/gophp"
import "github.com/boyxp/nova/cookie"
import "github.com/boyxp/nova/register"

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
	var data map[string]interface{}

	sess := read()
	if sess==nil {

		//生成sessionid，发送cookie
	}

	data[name] = value

//	byte, _ := gophp.Serialize(data)
//	string(jsonbyte)
//	think|
	return true
}

func ssid() string {
	ssid := cookie.Get("PHPSESSID")
	if ssid != "" {
		return ssid
	}

	ssid = cookie.Get("GOSESSID")
	if ssid != "" {
		return ssid
	}

	req := register.GetRequest()
	ip  := getIP(req)
	id  := getRoutineId()

    sum := md5.Sum([]byte(ip+":"+id))
    ssid = hex.EncodeToString(sum[:])

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

func getIP(req *http.Request) string {
    remoteAddr := req.RemoteAddr
    if ip := req.Header.Get("X-Real-IP"); ip != "" {
        remoteAddr = ip
    } else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
        remoteAddr = ip
    } else {
        remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
    }

    if remoteAddr == "::1" {
        remoteAddr = "127.0.0.1"
    }

    return remoteAddr
}

func getRoutineId() string {
    var buf [20]byte
    runtime.Stack(buf[:], false)
    for i:=10;i<20;i++ {
        if buf[i]==32 {
            return string(buf[10:i])
        }
    }

    return "1"
}
