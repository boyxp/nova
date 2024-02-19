package database

import "os"
import "log"
import "sync"
import "unicode"
import "runtime"
import "strings"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var cache sync.Map

func Register(tag string, dbname string, dsn string) {
	if tag=="database" && (dbname=="" || dsn=="") {
		log.Printf("\033[1;31;40m%s\033[0m\n",".env配置文件不存在或database.dbname和database.dsn未设置")
		os.Exit(1)
	}

	cache.Store("dsn."+tag, dsn)
	cache.Store("dbname."+tag, dbname)
}

func Dsn(tag string) string {
	value, ok := cache.Load("dsn."+tag);
    if !ok {
    	_dsn := os.Getenv(tag+".dsn")
    	if _dsn!="" {
    		cache.Store("dsn."+tag, _dsn)
    		return _dsn
    	} else {
        	panic("dsn配置不存在:"+tag)
    	}
    }

    dsn, _ := value.(string)

	return dsn
}

func Dbname(tag string) string {
	value, ok := cache.Load("dbname."+tag)
	if !ok {
		_dbname := os.Getenv(tag+".dbname")
    	if _dbname!="" {
    		cache.Store("dbname."+tag, _dbname)
    		return _dbname
    	} else {
        	panic("dbname配置不存在:"+tag)
    	}
	}

	dbname, _ := value.(string)

	return dbname
}

func NewOrm(tag ...string) *Orm {
	_, file, _, _ := runtime.Caller(1)
	path  := strings.Split(file, "/")
	model := strings.Replace(path[len(path)-1], ".go", "", 1)

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
	O  = O.Init(dbtag, table)

	return O
}

func Open(tag string) *sql.DB {
	dsn     := Dsn(tag)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	return db
}
