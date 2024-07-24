package database

import "os"
import "log"
import "strings"
import "strconv"
import "database/sql"

type Orm struct {
	db *sql.DB
	dbname string
	dbtag string
	table string
	primary string
	scheme map[string]string
	allFields []string
	total int

	selectFields string
	selectConds []string
	selectParams []interface{}
	selectPage int
	selectLimit int
	selectOrder []string
	selectGroup []string
	selectHaving string

	debug string
}

func (O *Orm) Init(dbtag string, table string) *Orm {
	O.dbtag  = dbtag
	O.dbname = Dbname(dbtag)
	O.table  = table
	O.debug  = os.Getenv("debug")
	O.db     = Open(dbtag)
	O.total  = -1

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
			maps, ok := conds[0].(map[string]interface{})
			if ok {
				for _field,_criteria := range maps {
						_string, ok := _criteria.(string)
						if ok {
							O.Where(_field, _string)
							continue
						}

						_set, ok := _criteria.([]interface{})
						if ok {
							_tmp := []interface{}{_field}
							_tmp  = append(_tmp, _set...)
							O.Where(_tmp...)
							continue
						}
				}
				return O
			} else {
				panic("第一个参数应为string类型或map[string]interface{}类型")
			}
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
										criteria = append(criteria, "_in_query_placeholder_")
										//panic("查询条件应为[]string类型,且至少存在一个元素")
									}

									placeholders := []string{}
									for _,v := range criteria {
										placeholders   = append(placeholders, "?")
										O.selectParams = append(O.selectParams, v)
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
		check1 := strings.Contains(" "+O.selectFields+" ", " "+field+" ")
		check2 := strings.Contains(" "+O.selectFields+",", " "+field+",")
		if !ok && !check1 && !check2 {
			panic(field+":聚合字段应为一般字段或聚合别名")
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

	_, ok  := O.scheme[field]
	check1 := strings.Contains(" "+O.selectFields+" ", " "+field+" ")
	check2 := strings.Contains(" "+O.selectFields+",", " "+field+",")
	if !ok && !check1 && !check2 {
		panic(field+":排序应为字段或聚合的别名")
	}

	O.selectOrder = append(O.selectOrder, field+" "+sort)

	return O
}

func (O *Orm) Page(page int) *Orm {
	if page < 1 {
		page = 1
	}

	O.selectPage = page

	return O
}

func (O *Orm) Limit(limit int) *Orm {
	if limit < 1 {
		limit = 20
	}

	O.selectLimit = limit

	return O
}

func (O *Orm) Result() *Result {
	list := O.Select()

	fields  := map[string]string{}
	for _,v := range O.selectColumns() {
		idx := strings.LastIndex(v, " ")
		v    = v[idx+1:]
		fields[v] = v
	}

	return &Result {
		Page   : O.selectPage,
		Limit  : O.selectLimit,
		orm    : O,
		list   : list,
		fields : fields,
	}
}

func (O *Orm) Select() []map[string]string {
	var result []map[string]string

	stmt := O.selectStmt()

	if O.debug=="yes" {
		log.Println("SQL:\t"+stmt)
		log.Println("PARAMS:\t",O.selectParams)
	}

	rows, err := O.db.Query(stmt, O.selectParams...)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

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

			key := columns[i]
			pos := strings.LastIndex(key, " ")
			if pos > -1 {
				key = key[pos+1:]
			}

			record[key] = value
		}

		result = append(result, record)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	return result
}

func (O *Orm) Find(primary ...string) map[string]string {
	var result map[string]string 

	if len(primary)>0 {
		O.Where(primary[0])
	}

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

func (O *Orm) Columns(fields ...string) map[string]string {
	var key   string
	var value string

	if len(fields)==0 {
		panic("参数不可为空")

	} else if(len(fields)==1) {
		key   = O.primary
		value = fields[0]

	} else {
		key   = fields[1]
		value = fields[0]
	}

	selectFields := O.selectFields

	_, ok_key := O.scheme[key]
	check_key1 := strings.Contains(" "+O.selectFields+" ", " "+key+" ")
	check_key2 := strings.Contains(" "+O.selectFields+",", " "+key+",")
	if !ok_key && !check_key1 && !check_key2 {
		panic(key+":应为字段或聚合的别名")
	}

	_, ok_value := O.scheme[value]
	check_value1 := strings.Contains(" "+O.selectFields+" ", " "+value+" ")
	check_value2 := strings.Contains(" "+O.selectFields+",", " "+value+",")
	if !ok_value && !check_value1 && !check_value2 {
		panic(value+":应为字段或聚合的别名")
	}

	if len(O.selectFields)==0 {
		O.selectFields = key+","+value
	}

	if strings.Contains(O.selectFields, "(") && len(O.selectGroup)==0 {
		panic("查询了聚合数据但未设置聚合字段")
	}

	list := O.Select()

	O.selectFields = selectFields


	var result = map[string]string{}
	if len(list)>0 {
		for _, v := range list {
			result[v[key]] = v[value]
		}
	}

	return result
}

func (O *Orm) Map(fields ...string) map[string]map[string]string {
	var key string

	if(len(fields)==0) {
		key   = O.primary
	} else {
		key   = fields[0]
	}

	selectFields := O.selectFields

	_, ok_key := O.scheme[key]
	check_key1 := strings.Contains(" "+O.selectFields+" ", " "+key+" ")
	check_key2 := strings.Contains(" "+O.selectFields+",", " "+key+",")
	if !ok_key && !check_key1 && !check_key2 {
		panic(key+":应为字段或聚合的别名")
	}

	if strings.Contains(O.selectFields, "(") && len(O.selectGroup)==0 {
		panic("查询了聚合数据但未设置聚合字段")
	}

	list := O.Select()

	O.selectFields = selectFields


	var result = map[string]map[string]string{}
	if len(list)>0 {
		_, ok := list[0][key]
		if !ok {
			panic(key+":key必须包含在读取字段里")
		}

		for _, v := range list {
			result[v[key]] = v
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

	if len(O.selectGroup)==0 {
		selectFields  := O.selectFields
		O.selectFields = "count(*) as aggs_count"

		total := O.Value("aggs_count")

		O.selectFields = selectFields

		conv,err := strconv.Atoi(total)
		if err == nil {
			result = conv
		}

	} else {
		last    := ""
		element := []string{}
		group   := ","+strings.Join(O.selectGroup, ",")+","
		fields  := strings.Split(O.selectFields, ",")

		for _, field := range fields {
			field = strings.TrimSpace(field)
			_, ok := O.scheme[field]
			if ok {
				element = append(element, field)
				continue
			}

			if strings.Contains(field, "(") && !strings.Contains(field, ")") {
				last = field
				continue
			}

			idx := strings.LastIndex(field, " ")
			if idx>-1 && strings.Contains(group, ","+field[idx+1:]+",") {
				end    := strings.LastIndex(field, ")")
				element = append(element, last+","+field[0:end+1])
				last    = ""
				continue
			}
		}




		selectFields  := O.selectFields
		selectGroup   := O.selectGroup
		selectOrder   := O.selectOrder

		O.selectGroup  = []string{}
		O.selectOrder  = []string{}
		O.selectFields = "count(distinct "+strings.Join(element, ",")+") as aggs_count"

		total := O.Value("aggs_count")

		O.selectFields = selectFields
		O.selectGroup  = selectGroup
		O.selectOrder  = selectOrder

		conv,err := strconv.Atoi(total)
		if err == nil {
			result = conv
		}
	}

	O.total = result

	return result
}

func (O *Orm) Total() int {
	if O.total<0 {
		O.Count()
	}

	return O.total
}

func (O *Orm) TotalPage() int {
	if O.total<0 {
		O.Count()
	}

	if O.selectLimit==0 {
		O.selectLimit = 20
	}

	var total_page int
	if O.total==0 {
		total_page = 0

	} else if O.total%O.selectLimit==0 {
		total_page = O.total/O.selectLimit

	} else {
		total_page = int(O.total/O.selectLimit)+1
	}

	return total_page
}

func (O *Orm) Exist(primary string) bool {
	res := O.Where(primary).Find()
	return res!=nil
}

func (O *Orm) Relate(list *[]map[string]string, fields string) {
	ids := []string{"0"}
	for _, v := range *list {
		if _id, ok := v[O.primary];ok {
			ids = append(ids, _id)
		} else {
			panic("记录列表缺少主键字段:"+O.primary)
		}
	}


	_empty  := map[string]string{}
	_fields := strings.Split(fields, ",")
	for _, field := range _fields {
		_, ok := O.scheme[field]
		if !ok {
			panic(field+":字段不存在")
		}

		_empty[field] = ""
	}


	_result := O.Field(O.primary+","+fields).Where(O.primary, "in", ids).Limit(10000).Select()
	_temp   := map[string]map[string]string{}
	for _, v := range _result {
		_temp[v[O.primary]] = v
	}


	merge := map[string]string{}
	for _, v := range *list {
		key := v[O.primary]
		if r, ok := _temp[key];ok {
			merge = r
		} else {
			merge = _empty
		}

		for _k, _v := range merge {
			v[_k] = _v
		}
	}
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

	if strings.Contains(O.selectFields, "aggs_count") {
		return []string{"aggs_count"}
	}

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

	var i int  = 1
	var length = len(data)
	for f,v := range data {

		if _, ok := O.scheme[f];!ok {
			panic(f+":字段不存在")
		}

		sql.WriteString(f+"=? ")
		if i < length {
			sql.WriteString(",")
		}
		params = append(params, v)
		i = i+1
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
	if O.debug=="yes" {
		log.Println("SQL:\t"+_sql)
	}

	stmt, err := O.db.Prepare(_sql)
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
		rows, err := O.db.Query("describe "+table)
		if err != nil {
    		panic(err.Error())
		}

		defer rows.Close()

		for rows.Next() {
			var rowField, rowType, rowNull, rowKey, rowExtra string
			var rowDefault sql.RawBytes

			if err := rows.Scan(&rowField, &rowType, &rowNull, &rowKey, &rowDefault, &rowExtra); err != nil {
        		panic(err.Error())
    		}

			scheme[rowField] = rowType

			if rowKey=="PRI" {
				primary = rowField
			}

			allFields = append(allFields, rowField)
		}

		cache.Store("primary."+O.dbname+"."+table, primary)
		cache.Store("allFields."+O.dbname+"."+table, allFields)
		cache.Store("scheme."+O.dbname+"."+table, scheme)

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
