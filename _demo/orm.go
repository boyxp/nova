package main

import (
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	o := Orm{table:"monitor",primary:"id", scheme:map[string]string{"test":"string","create_at":"timestamp"}}
	//o.Insert(map[string]interface{}{"test":"111"})
	o.Field("test,create_at").
	Where("11").
	Where("test","22").
	Where("test", ">=", "33").
	Where("test","in",[]string{"44","55","66"}).
	Where("test","is","null").
	Where("test","BETWEEN", []string{"77","88"})
}

type Orm struct {
	db *sql.DB
	table string
	fields string
	primary string
	scheme map[string]string
	conds []string
	params []interface{}
}

func (O *Orm) Insert(data map[string]interface{}) int64 {
	O.open()
	defer O.close()

	fields       := []string{}
	placeholders := []string{}
	values       := []interface{}{}

	for k,v := range data {
		fields       = append(fields, k)
		placeholders = append(placeholders, "?")
		values       = append(values, v)
	}

	stmt, err := O.db.Prepare("INSERT INTO "+O.table+" ("+strings.Join(fields, ",")+") VALUES("+strings.Join(placeholders, ",")+")")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec(values...)
	if err != nil {
		panic(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return id
}

func (O *Orm) open(){
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/dev_xinzhanghu")
	if err != nil {
		panic(err.Error())
	}

	O.db = db
}

func (O *Orm) close() {
	O.db.Close()
}

func (O *Orm) Delete() {

}

func (O *Orm) Update() {

}

func (O *Orm) Select() {

}

func (O *Orm) Find() {

}

func (O *Orm) Value() {

}

func (O *Orm) Execute() {

}

func (O *Orm) Sum() {

}

func (O *Orm) Count() {

}

func (O *Orm) Field(fields string) *Orm {
	O.fields = fields
	return O
}

func (O *Orm) Where(conds ...interface{}) *Orm {
	args_len := len(conds)
	switch args_len {
		case 1 :
				id, ok := conds[0].(string)
				if !ok {
					panic("查询条件应为string类型")
				}

				O.conds  = append(O.conds, O.primary+"=?")
				O.params = append(O.params, id)
		case 2 :
				field, ok := conds[0].(string)
				if !ok {
					panic("查询字段应为string类型")
				}

				_, ok = O.scheme[field]
				if !ok {
					panic(field+":查询字段不存在")
				}

				criteria, ok := conds[1].(string)
				if !ok {
					panic("查询值应为string类型")
				}

				O.conds  = append(O.conds, field+"=?")
				O.params = append(O.params, criteria)
		case 3 :
				field, ok := conds[0].(string)
				if !ok {
					panic("查询字段应为string类型")
				}

				_, ok = O.scheme[field]
				if !ok {
					panic(field+":查询字段不存在")
				}

				opr, ok := conds[1].(string)
				if !ok {
					panic("运算符应为string类型")
				}

				opr = strings.ToTitle(opr)
				switch opr {
					case "!="     : fallthrough
					case ">"      : fallthrough
					case ">="     : fallthrough
					case "<"      : fallthrough
					case "<="     : 
									criteria, ok := conds[2].(string)
									if !ok {
										panic("查询条件应为string类型")
									}

									O.conds  = append(O.conds, field+opr+"?")
									O.params = append(O.params, criteria)

					case "IN"     : fallthrough
					case "NOT IN" :
									criteria, ok := conds[2].([]string)
									if !ok {
										panic("查询条件应为[]string类型")
									}

									if len(criteria)==0 {
										panic("查询条件应为[]string类型,且至少存在一个元素")
									}

									placeholders := []string{}
									for _,v := range criteria {
										placeholders = append(placeholders, "?")
										O.params     = append(O.params, v)
									}

									O.conds  = append(O.conds, field+" "+opr+"("+strings.Join(placeholders, ",")+")")

					case "IS"     : fallthrough
					case "IS NOT" :
									criteria, ok := conds[2].(string)
									if !ok {
										panic("查询条件应为string类型")
									}

									criteria = strings.ToTitle(criteria)
									if criteria!="NULL" {
										panic("查询条件只能为null")
									}

									O.conds  = append(O.conds, field+" "+opr+" "+criteria)

					case "BETWEEN":
									criteria, ok := conds[2].([]string)
									if !ok {
										panic("查询条件应为[]string类型")
									}

									if len(criteria)!=2 {
										panic("查询条件应为[]string类型,且必须2个元素")
									}

									O.conds  = append(O.conds, field+" "+opr+" ? AND ? ")
									for _,v := range criteria {
										O.params     = append(O.params, v)
									}
					case "LIKE"   :
					case "EXP"    :
				}
		default : panic("查询参数不应超过3个")
	}

	fmt.Println(O.conds)
	fmt.Println(O.params)

	return O
}

func (O *Orm) Page() {

}

func (O *Orm) Limit() {

}

func (O *Orm) Order() {

}

func (O *Orm) Group() {

}

func (O *Orm) Having() {

}
