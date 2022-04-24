package database

import "strings"
import "strconv"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Orm struct {
	dbname string
	dbtag string
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

func (O *Orm) Init(dbtag string, table string) *Orm {
	O.dbname = Dbname(dbtag)
	O.dbtag  = dbtag

	scheme, primary := O.getScheme(table)

	O.table        = table
	O.scheme       = scheme
	O.primary      = primary
	O.selectFields = ""
	O.selectConds  = []string{}
	O.selectParams = []interface{}{}
	O.selectLimit  = 20

	return O
}

func (O *Orm) Insert(data map[string]interface{}) int64 {
	db := Open(O.dbtag)
	defer db.Close()

	fields       := []string{}
	placeholders := []string{}
	values       := []interface{}{}

	for k,v := range data {
		fields       = append(fields, k)
		placeholders = append(placeholders, "?")
		values       = append(values, v)
	}

	stmt, err := db.Prepare("INSERT INTO "+O.table+" ("+strings.Join(fields, ",")+") VALUES("+strings.Join(placeholders, ",")+")")
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

func (O *Orm) Delete() int64 {
	db := Open(O.dbtag)
	defer db.Close()

	sql := O.deleteStmt()
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec(O.selectParams...)
	if err != nil {
		panic(err.Error())
	}

	ar, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	return ar
}

func (O *Orm) Update(data map[string]string) int64 {
	db := Open(O.dbtag)
	defer db.Close()

	if len(data)==0 {
		panic("没有更新字段")
	}

	sql, params := O.updateStmt(data)

	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec(params...)
	if err != nil {
		panic(err.Error())
	}

	ar, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	return ar
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
			if k==0 {
				continue
			}

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

func (O *Orm) Select() []map[string]string {
	var result []map[string]string

	db := Open(O.dbtag)
	defer db.Close()

	stmt := O.selectStmt()
	rows, err := db.Query(stmt, O.selectParams...)
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

		var record = map[string]string{}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			record[columns[i]] = value
		}

		result = append(result, record)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	return result
}

func (O *Orm) Find() map[string]string {
	var result map[string]string 

	selectPage   := O.selectPage
	selectLimit  := O.selectLimit
	O.selectPage  = 1
	O.selectLimit = 1
	list := O.Select()
	O.selectPage  = selectPage
	O.selectLimit = selectLimit

	if len(list)>0 {
		result = list[0]
	}

	return result
}

func (O *Orm) Value(field string) string {
	selectFields := O.selectFields

	_, ok := O.scheme[field]
	if ok {
		O.selectFields = field
	} else if strings.Contains(" "+O.selectFields+" ", " "+field+" ") {
	} else {
		panic(field+":取值字段应为普通字段或聚合别名")
	}

	row := O.Find()
	O.selectFields = selectFields
	value, ok := row[field]
	if ok {
		return value
	}

	return ""
}

func (O *Orm) Sum(field string) int {
	_, ok := O.scheme[field]
	if !ok {
		panic(field+":聚合字段不存在")
	}

	var result int
	selectFields  := O.selectFields
	O.selectFields = "sum("+field+") as aggs_sum"

	total := O.Value("aggs_sum")
	O.selectFields = selectFields
	conv,err := strconv.Atoi(total)
	if err == nil {
		result = conv
	}

	return result
}

func (O *Orm) Count() int {
	var result int
	selectFields  := O.selectFields
	O.selectFields = "count(*) as aggs_count"

	total := O.Value("aggs_count")
	O.selectFields = selectFields
	conv,err := strconv.Atoi(total)
	if err == nil {
		result = conv
	}

	return result
}

func (O *Orm) Execute() {

}




func (O *Orm) selectStmt() string {
	var sql strings.Builder

	sql.WriteString("SELECT ")
	if O.selectFields == "" {
		sql.WriteString("* ")
	} else {
		if strings.Contains(strings.ToTitle(O.selectFields), " COUNT(") && strings.Contains(O.selectFields, ",") && len(O.selectGroup)==0 {
			panic("缺少聚合字段")
		}

		if strings.Contains(strings.ToTitle(O.selectFields), " SUM(") && strings.Contains(O.selectFields, ",") && len(O.selectGroup)==0 {
			panic("缺少聚合字段")
		}

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
		if !strings.Contains(O.selectFields, ")") {
			panic("缺少聚合查询")
		}

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

func (O *Orm) deleteStmt() string {
	var sql strings.Builder

	sql.WriteString("DELETE FROM ")
	sql.WriteString(O.table)
	if len(O.selectConds)>0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(O.selectConds, " AND "))
		sql.WriteString(" ")
	}

	if O.selectLimit==0 {
		sql.WriteString("LIMIT 20")
	} else {
		sql.WriteString("LIMIT ")
		sql.WriteString(strconv.Itoa(O.selectLimit))
	}

	return sql.String()
}

func (O *Orm) updateStmt(data map[string]string) (string,[]interface{}) {
	var sql strings.Builder
	var params []interface{}

	sql.WriteString("UPDATE ")
	sql.WriteString(O.table)
	sql.WriteString(" SET ")

	for f,v := range data {
		sql.WriteString(f+"=? ")
		params = append(params, v)
	}
	params = append(params, O.selectParams...)

	if len(O.selectConds)>0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(O.selectConds, " AND "))
		sql.WriteString(" ")
	}

	if O.selectLimit==0 {
		sql.WriteString("LIMIT 20")
	} else {
		sql.WriteString("LIMIT ")
		sql.WriteString(strconv.Itoa(O.selectLimit))
	}

	return sql.String(), params
}

func (O *Orm) getScheme(table string) (map[string]string, string) {
	var primary string
	var scheme = map[string]string{}

	value, ok := cache.Load("scheme."+O.dbname+"."+table)
    if !ok {
		O.table  = "information_schema.columns"
		O.scheme = map[string]string{"TABLE_SCHEMA":"","TABLE_NAME":"","COLUMN_NAME":"","IS_NULLABLE":"","COLUMN_DEFAULT":"","COLUMN_KEY":""}
		columns := O.Field("COLUMN_NAME,IS_NULLABLE,COLUMN_DEFAULT,COLUMN_KEY").
				Where("TABLE_SCHEMA",O.dbname).
				Where("TABLE_NAME", table).
				Limit(200).
				Select()

		for _,r := range columns {
			if r["IS_NULLABLE"]=="NO" && r["COLUMN_DEFAULT"]=="NULL" {
				scheme[r["COLUMN_NAME"]] = "NOTNULL"
			} else {
				scheme[r["COLUMN_NAME"]] = "NULL"
			}

			if r["COLUMN_KEY"]=="PRI" {
				primary = r["COLUMN_NAME"]
			}
		}

		cache.Store("scheme."+O.dbname+"."+table, scheme)
		cache.Store("primary."+O.dbname+"."+table, primary)
	} else {
		scheme, _  = value.(map[string]string)
		value, _  := cache.Load("primary."+O.dbname+"."+table)
		primary, _ = value.(string)
	}

	return scheme, primary
}

