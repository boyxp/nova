package database

import "strings"
import "strconv"
import "database/sql"

type Orm struct {
	dbname string
	dbtag string
	table string
	primary string
	scheme map[string]string
	allFields []string

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
	O.dbtag  = dbtag
	O.dbname = Dbname(dbtag)
	O.table  = table

	O.initScheme(table)

	return O
}

func (O *Orm) Insert(data map[string]string) int64 {
	fields       := []string{}
	placeholders := []string{}
	values       := []interface{}{}

	for k,v := range data {
		if _, ok := O.scheme[k];!ok {
			panic(k+":字段不存在")
		}

		fields       = append(fields, k)
		placeholders = append(placeholders, "?")
		values       = append(values, v)
	}

	_sql := "INSERT INTO "+O.table+" ("+strings.Join(fields, ",")+") VALUES("+strings.Join(placeholders, ",")+")"
	res  := O.execute(_sql, values)

	id, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return id
}

func (O *Orm) Delete() int64 {
	_sql := O.deleteStmt()
	res  := O.execute(_sql, O.selectParams)

	ar, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	return ar
}

func (O *Orm) Update(data map[string]string) int64 {
	if len(data)==0 {
		panic("没有更新字段")
	}

	_sql, params := O.updateStmt(data)

	res := O.execute(_sql, params)

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

func (O *Orm) Group(fields ...string) *Orm {
	for _, field := range fields {
		_, ok := O.scheme[field]
		if !ok {
			panic(field+":聚合字段不存在")
		}

		O.selectGroup = append(O.selectGroup, field)
	}

	return O
}

func (O *Orm) Having(field string, opr string, criteria int) *Orm {
	if len(O.selectGroup)==0 {
		panic("没有聚合字段")
	}

	_, ok := O.scheme[field]
	check1 := strings.Contains(" "+O.selectFields+" ", " "+field+" ")
	check2 := strings.Contains(" "+O.selectFields+",", " "+field+",")
	if !ok && !check1 && !check2 {
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

		var value string
		var record = map[string]string{}
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

func (O *Orm) Values(field string) []string {
	var result []string

	list := O.Select()

	if len(list)>0 {
		for _, v := range list {
			result = append(result, v[field])
		}
	}

	return result
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

func (O *Orm) Exist(primary string) bool {
	res := O.Where(primary).Find()
	return res!=nil
}


func (O *Orm) selectStmt() string {
	var sql strings.Builder

	sql.WriteString("SELECT ")
	if O.selectFields == "" {
		sql.WriteString(strings.Join(O.allFields, ",")+" ")
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
		return O.allFields
	}

	if !strings.Contains(O.selectFields," ") && !strings.Contains(O.selectFields, "(") {
		return strings.Split(O.selectFields, ",")

	} else {
		if strings.Contains(O.selectFields, "(")==false {
			for _,v := range strings.Split(O.selectFields,",") {
				columns = append(columns, strings.TrimSpace(v))
			}

		} else {
			var skip = false
			for _,v := range strings.Split(O.selectFields,",") {
				if strings.Contains(v, "(") {
					skip = true
				}

				if strings.Contains(v, ")") {
					tmp := strings.Split(strings.TrimSpace(v)," ")
					columns = append(columns, tmp[len(tmp)-1])
					skip = false
					continue
				}

				if skip {
					continue
				}

				columns = append(columns, strings.TrimSpace(v))
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
		if _, ok := O.scheme[f];!ok {
			panic(f+":字段不存在")
		}

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

func (O *Orm) execute(_sql string, values []interface{}) sql.Result {
	db := Open(O.dbtag)
	defer db.Close()

	stmt, err := db.Prepare(_sql)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec(values...)
	if err != nil {
		panic(err.Error())
	}

	return res
}

func (O *Orm) initScheme(table string) {
	var primary string
	var allFields []string
	var scheme = map[string]string{}

	value, ok := cache.Load("scheme."+O.dbname+"."+table)
    if !ok {
    	db := Open(O.dbtag)
		defer db.Close()

		rows, err := db.Query("describe "+table)
		if err != nil {
    		panic(err.Error())
		}

		for rows.Next() {
			var rowField, rowType, rowNull, rowKey, rowExtra string
			var rowDefault sql.RawBytes

			if err := rows.Scan(&rowField, &rowType, &rowNull, &rowKey, &rowDefault, &rowExtra); err != nil {
        		panic(err.Error())
    		}

			if rowNull=="NO" && rowDefault==nil {
				scheme[rowField] = "NOTNULL"
			} else {
				scheme[rowField] = "NULL"
			}

			if rowKey=="PRI" {
				primary = rowField
			}

			allFields = append(allFields, rowField)
		}

		cache.Store("scheme."+O.dbname+"."+table, scheme)
		cache.Store("primary."+O.dbname+"."+table, primary)
		cache.Store("allFields."+O.dbname+"."+table, allFields)

	} else {
		scheme, _    = value.(map[string]string)
		value1, _   := cache.Load("primary."+O.dbname+"."+table)
		primary, _   = value1.(string)
		value2, _   := cache.Load("allFields."+O.dbname+"."+table)
		allFields, _ = value2.([]string)
	}

	O.primary   = primary
	O.scheme    = scheme
	O.allFields = allFields
}
