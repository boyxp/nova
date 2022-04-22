package main

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	o := Orm{table:"payment",primary:"payment_id", scheme:map[string]string{"payment_id":"int","company_id":"int","batch_id":"int","receive_name":"string","create_at":"timestamp","tax_mobile":"string"}}

	//o.Where("1").Select()
	o.Field("company_id  ,  receive_name,tax_mobile").Where("1").Select()
	//o.Field("company_id,count(*) as total").Where("1").Select()

	//o.Insert(map[string]interface{}{"test":"111"})
	//o.Field("company_id,count(*) as total").
	//Where("1").
	//Where("company_id","22").
	//Where("company_id", ">=", "33").
	//Where("company_id","in",[]string{"44","55","66"}).
	//Where("company_id","is","null").
	//Where("company_id","BETWEEN", []string{"77","88"}).
	//Where("company_id in (?,?) and test>?", "99","100","101").
	//Where("receive_name","like","abc").
	//Page(1).
	//Limit(5).
	//Order("total","asc").
	//Group("company_id").
	//Group("batch_id").
	//Having("total",">",1).
	//Select()
}

type Orm struct {
	db *sql.DB
	table string
	primary string
	scheme map[string]string

	selectFields string
	selectConds []string
	selectParams []interface{}
	selectPage int
	selectLimit int
	selectOrder []string
	selectGroup []string
	selectHaving string
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

func (O *Orm) Execute() {

}

func (O *Orm) Sum() {

}

func (O *Orm) Count() {

}

func (O *Orm) Field(fields string) *Orm {
	O.selectFields = fields
	return O
}

func (O *Orm) Where(conds ...interface{}) *Orm {
	args_len := len(conds)
	if args_len < 1 {
		panic("查询参数不应为空")
	}

	field, ok := conds[0].(string)
	if !ok {
		panic("第一个参数应为string类型")
	}

	if placeholder_count := strings.Count(field, "?"); placeholder_count > 0 {
		if placeholder_count != args_len -1 {
			panic("查询占位符和参数数量不匹配")
		}

		for k,v := range conds {
			param, ok := v.(string)
			if !ok {
				panic("第"+strconv.Itoa(k)+"个参数应为string类型")
			}

			O.selectParams = append(O.selectParams, param)
		}

		O.selectConds  = append(O.selectConds, field)

		return O
	}


	switch args_len {
		case 1 :
				O.selectConds  = append(O.selectConds, O.primary+"=?")
				O.selectParams = append(O.selectParams, field)
		case 2 :
				_, ok = O.scheme[field]
				if !ok {
					panic(field+":查询字段不存在")
				}

				criteria, ok := conds[1].(string)
				if !ok {
					panic("查询值应为string类型")
				}

				O.selectConds  = append(O.selectConds, field+"=?")
				O.selectParams = append(O.selectParams, criteria)
		case 3 :
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
					case "<="     : fallthrough
					case "LIKE"   :
									criteria, ok := conds[2].(string)
									if !ok {
										panic("查询条件应为string类型")
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+" ?")
									O.selectParams = append(O.selectParams, criteria)

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
										O.selectParams     = append(O.selectParams, v)
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+"("+strings.Join(placeholders, ",")+")")

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

									O.selectConds  = append(O.selectConds, field+" "+opr+" "+criteria)

					case "BETWEEN":
									criteria, ok := conds[2].([]string)
									if !ok {
										panic("查询条件应为[]string类型")
									}

									if len(criteria)!=2 {
										panic("查询条件应为[]string类型,且必须2个元素")
									}

									O.selectConds  = append(O.selectConds, field+" "+opr+" ? AND ? ")
									for _,v := range criteria {
										O.selectParams     = append(O.selectParams, v)
									}

					default        :
									panic("不支持的操作类型:"+opr)
				}
		default : panic("查询参数不应超过3个")
	}

	fmt.Println(O.selectConds)
	fmt.Println(O.selectParams)

	return O
}

func (O *Orm) Page(page int) *Orm {
	if page < 1 {
		panic("页码不应小于1")
	}

	O.selectPage = page
	return O
}

func (O *Orm) Limit(limit int) *Orm {
	if limit < 1 {
		panic("每页条数不应小于1")
	}

	O.selectLimit = limit
	return O
}

func (O *Orm) Order(field string, sort string) *Orm {
	sort = strings.ToTitle(sort)
	if sort!="DESC" && sort!="ASC" {
		panic("排序类型只能是asc或desc")
	}

	_, ok := O.scheme[field]
	check := strings.Contains(" "+O.selectFields+" ", " "+field+" ")
	if !ok && !check {
		panic(field+":排序应为字段或聚合的别名")
	}

	O.selectOrder = append(O.selectOrder, field+" "+sort)

	return O
}

func (O *Orm) Group(field string) *Orm {
	_, ok := O.scheme[field]
	if !ok {
		panic(field+":聚合字段不存在")
	}

	O.selectGroup = append(O.selectGroup, field)

	return O
}

func (O *Orm) Having(field string, opr string, criteria int) *Orm {
	if len(O.selectGroup)==0 {
		panic("没有聚合字段")
	}

	_, ok := O.scheme[field]
	check := strings.Contains(" "+O.selectFields+" ", " "+field+" ")
	if ok || !check {
		panic(field+":过滤字段应为聚合别名")
	}

	if !strings.Contains(",=,!=,>,>=,<,<=,", ","+opr+",") {
		panic("不支持的过滤操作符号:"+opr)
	}

	O.selectHaving = field+" "+opr+" "+strconv.Itoa(criteria)
	fmt.Println(O.selectHaving)

	return O
}

func (O *Orm) Select() {
	O.open()
	defer O.close()

	stmt := O.selectStmt()
	

	rows, err := O.db.Query(stmt, O.selectParams...)
	if err != nil {
		panic(err.Error())
	}

	columns  := O.selectColumns()
	values   := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			fmt.Println(columns[i], "=", value)
		}
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	fmt.Println(stmt)
}
func (O *Orm) selectStmt() string {
	var sql strings.Builder

	sql.WriteString("SELECT ")
	if O.selectFields == "" {
		sql.WriteString("* ")
	} else {
		sql.WriteString(O.selectFields)
		sql.WriteString(" ")
	}

	sql.WriteString("FROM ")
	sql.WriteString(O.table)
	sql.WriteString(" ")

	if len(O.selectConds)>0 {
		sql.WriteString("WHERE ")
		sql.WriteString(strings.Join(O.selectConds, " AND "))
		sql.WriteString(" ")
	}

	if len(O.selectGroup)>0 {
		sql.WriteString("GROUP BY ")
		sql.WriteString(strings.Join(O.selectGroup, ","))
		sql.WriteString(" ")
	}

	if O.selectHaving!="" {
		sql.WriteString("HAVING ")
		sql.WriteString(O.selectHaving)
		sql.WriteString(" ")
	}

	if len(O.selectOrder)>0 {
		sql.WriteString("ORDER BY ")
		sql.WriteString(strings.Join(O.selectOrder,","))
		sql.WriteString(" ")
	}


	var offset int
	if O.selectLimit==0 {
		O.selectLimit = 20
	}

	if O.selectPage<=1 {
		O.selectPage = 1
		offset = 0
	} else {
		offset = O.selectPage*O.selectLimit-O.selectLimit
	}

	sql.WriteString("LIMIT ")
	sql.WriteString(strconv.Itoa(offset))
	sql.WriteString(",")
	sql.WriteString(strconv.Itoa(O.selectLimit))
	sql.WriteString("  ")

	return sql.String()
}
func (O *Orm) selectColumns() []string {
	var columns []string
	if O.selectFields == "" || O.selectFields=="*" {
		for k,_ := range O.scheme {
			columns = append(columns, k)
		}
		return columns
	}

	if !strings.Contains(O.selectFields," ") {
		return strings.Split(O.selectFields,",")
	} else {
		for _,v := range strings.Split(O.selectFields,",") {
			if !strings.Contains(v, " ") {
				columns = append(columns, v)
			} else {
				tmp := strings.Split(strings.TrimSpace(v)," ")
				columns = append(columns, tmp[len(tmp)-1])
			}
		}
	}

	return columns
}

func (O *Orm) Find() {

}

func (O *Orm) Value() {

}

