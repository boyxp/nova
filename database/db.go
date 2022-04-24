package database

import "sync"
import "unicode"
import "runtime"
import "strings"

var cache sync.Map
func Register(tag string, dbname string, dsn string) {
	cache.Store("dsn."+tag, dsn)
	cache.Store("dbname."+tag, dbname)
}

func Dsn(tag string) string {
	value, ok := cache.Load("dsn."+tag)
    if !ok {
        panic("dsn配置不存在:"+tag)
    }

    dsn, _ := value.(string)

	return dsn
}

func Dbname(tag string) string {
	value, ok := cache.Load("dbname."+tag)
	if !ok {
		panic("dbname配置不存在:"+tag)
	}

	dbname, _ := value.(string)

	return dbname
}

func NewOrm(tag ...string) *Orm {
	_, file, _, _ := runtime.Caller(1)
	path  := strings.Split(file, "/")
	model := strings.Trim(path[len(path)-1], ".go")

	var dbtag string
	if len(tag)==0 {
		dbtag = "database"
	} else {
		dbtag = tag[0]
	}

	var table string
	value, ok := cache.Load(dbtag+"."+model)
    if !ok {
		var buf strings.Builder
		for i, l := range model {
			if unicode.IsUpper(l) {
				if i != 0 {
  					buf.WriteString("_")
				}
 				buf.WriteString(string(unicode.ToLower(l)))
			} else {
 				buf.WriteString(string(l))
			}
		}

		table = buf.String()
		cache.Store(dbtag+"."+model, table)
    } else {
    	table, _ = value.(string)
    }

	O := &Orm{}
	O = O.Init(dbtag, table)

	return O
}
