package session

import "os"
import "net"
import "runtime"
import "strconv"
import "strings"
import "net/http"
import "io/ioutil"
import "crypto/md5"
import "encoding/hex"
import "github.com/leeqvip/gophp/serialize"
import "github.com/boyxp/nova/cookie"
import "github.com/boyxp/nova/register"
import "github.com/boyxp/nova/cache"

func All() map[string]string {
	result := map[string]string{}
	sess, _, _ := read("r")

	for k, v := range sess {
		if k=="_t" {
			continue
		}

		switch v.(type) {
			case string :
					value, ok := v.(string)
					if !ok {
						result[k] = ""
					}
					result[k] = value

			case int   :
				   value, ok := v.(int)
				   if !ok {
						result[k] = ""
				   }
				   result[k] = strconv.Itoa(value)
		}
	}

	return result
}

func Get(name string) string {
	sess, _, _ := read("r")
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
	var data map[string]any

	sess, ssid, path := read("w")

	if sess==nil {
		data = make(map[string]any)
	} else {
		data = sess
	}

	cookie.Set("PHPSESSID", ssid)

	data[name] = value

	_byte, _ := serialize.Marshal(data)
	content := "think|"+string(_byte)

	real := path+"/sess_"+ssid
	err  := ioutil.WriteFile(real, []byte(content), 0666)

	return err == nil
}

func Id() string {
	return getSsid()
}

func getSsid() string {
	var ssid string

	req := register.GetRequest()

	c1, err1 := req.Cookie("PHPSESSID")
	if err1 == nil {
		return c1.Value
	}

	c2, err2 := req.Cookie("GOSESSID")
	if err2 == nil {
		return c2.Value
	}
	
	ssid = req.Header.Get("X-SESSID")
	if ssid != "" {
		return ssid
	}

	ip  := getIP(req)
	id  := getRoutineId()

    sum := md5.Sum([]byte(ip+":"+id))
    ssid = hex.EncodeToString(sum[:])

    return ssid
}

func getPath() string {
	_path := os.Getenv("session.path")
	if _path != "" {
		return _path
	}

	return "/var/lib/php/session/"
}

func read(mode string) (map[string]any, string, string) {
	var empty map[string]any

	ssid := getSsid()
	path := getPath()
	real := path+"/sess_"+ssid

	//cache
	if mode=="r" {
		info, er := os.Stat(real)
	    if er == nil {
	        memo := cache.Memory{}.Get("sess_"+ssid)
	        if memo!=nil {
	        	modTime := info.ModTime()
	        	sess    := memo.(map[string]any)
	        	if modTime==sess["_t"] {
	        		return sess, ssid, path
	        	}
	        }
	    }
	}

	file, err := os.Open(real)
	defer file.Close()
   	if err != nil {
        return empty, ssid, path
   	}

   content, err := ioutil.ReadAll(file)
   str := string(content)
   if len(str)==0 {
   	    return empty, ssid, path
   }

   if strings.Contains(str, "|") {
   		str = str[strings.Index(str, "|")+1:]
   }

	data, _  := serialize.UnMarshal([]byte(str))
	sess, ok := data.(map[string]any)
	if !ok {
		return empty, ssid, path
	}

	//cache
	if mode=="r" {
		info, er := os.Stat(real)
    	if er == nil {
        	sess["_t"] = info.ModTime()
        	cache.Memory{}.Set("sess_"+ssid, sess, 3600)
    	}
	}

	return sess, ssid, path
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
