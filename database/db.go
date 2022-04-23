package database

import "sync"

var dsnMap sync.Map
var dbnameMap sync.Map
func Register(tag string, dbname string, dsn string) {
	dsnMap.Store(tag, dsn)
	dbnameMap.Store(tag, dbname)
}

func Dsn(tag string) string {
	value, ok := dsnMap.Load(tag)
    if !ok {
        panic("dsn配置不存在:"+tag)
    }

    dsn, _ := value.(string)

	return dsn
}

func Dbname(tag string) string {
	value, ok := dbnameMap.Load(tag)
	if !ok {
		panic("dbname配置不存在:"+tag)
	}

	dbname, _ := value.(string)

	return dbname
}
