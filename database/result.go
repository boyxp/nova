package database

type Result struct {
	Page int
	Limit int
	orm *Orm
	list []map[string]string
	fields map[string]string
}

//返回查询结果切片
func (R *Result) List() []map[string]string {
	return R.list
}

//返回指定列切片
func (R *Result) Values(field string) []string {
	if _,ok := R.fields[field];!ok {
		panic("结果集中不存在指定字段或别名："+field)
	}

	res := []string{}

	for _,v := range R.list {
		res = append(res, v[field])
	}

	return res
}

//返回指定列做key指定列做value的map
func (R *Result) Columns(fields ...string) map[string]string {
	var key   string
	var value string

	if len(fields)==0 {
		panic("参数不可为空")

	} else if(len(fields)==1) {
		key   = R.orm.primary
		value = fields[0]

	} else {
		key   = fields[1]
		value = fields[0]
	}

	if _,ok_key := R.fields[key];!ok_key {
		panic("结果集中不存在指定字段或别名（未指定key时默认主键为key）："+key)
	}

	if _,ok_value := R.fields[value];!ok_value {
		panic("结果集中不存在指定字段或别名："+value)
	}

	res := map[string]string{}

	if len(R.list)>0 {
		for _, v := range R.list {
			res[v[key]] = v[value]
		}
	}

	return res
}

//返回指定列为key记录为值的map
func (R *Result) Map(fields ...string) map[string]map[string]string {
	var key string

	if(len(fields)==0) {
		key = R.orm.primary
	} else {
		key = fields[0]
	}

	if _,ok_key := R.fields[key];!ok_key {
		panic("结果集中不存在指定字段或别名（未指定key时默认主键为key）："+key)
	}

	res := map[string]map[string]string{}

	if len(R.list)>0 {
		for _, v := range R.list {
			res[v[key]] = v
		}
	}

	return res
}

//返回指定列为key指定列为值的切片
func (R *Result) MapList(key string, value string) map[string][]string {
	if _,ok_key := R.fields[key];!ok_key {
		panic("结果集中不存在指定字段或别名："+key)
	}

	if _,ok_value := R.fields[value];!ok_value {
		panic("结果集中不存在指定字段或别名："+value)
	}

	res := map[string][]string{}

	if len(R.list)>0 {
		for _, v := range R.list {
			if _, ok := res[v[key]];ok {
				res[v[key]] = append(res[v[key]], v[value])
			} else {
				res[v[key]] = []string{v[value]}
			}
		}
	}

	return res
}

//返回记录总条数
func (R *Result) Total() int {
	return R.orm.Total()
}

//返回记录总页数
func (R *Result) TotalPage() int {
	return R.orm.TotalPage()
}

//返回结果集是否为空
func (R *Result) Empty() bool {
	return len(R.list)==0
}

//找到两个结果集的唯一交集key，合并结果集
func (R *Result) Merge(M *Result) []map[string]string {
	//查找结果集字段交集
	fields := M.Fields()
	keys   := []string{}
	for k, _ := range fields {
		if _, ok := R.fields[k];ok {
			keys = append(keys, k)
		}
	}

	if len(keys)==0 {
		panic("两个结果集没有字段交集")
	}

	if len(keys)>1 {
		panic("两个结果集只能一个字段交集")
	}

	//按交集key转map
	keymap := M.Map(keys[0])

	//遍历合并，更新结果
	ik := keys[0]
	for rk, rv := range R.list {
		if _, ok := keymap[rv[ik]];ok {
			for mk, mv := range keymap[rv[ik]] {
				R.list[rk][mk] = mv
			}
		} else {
			for fk, _ := range fields {
				R.list[rk][fk] = ""
			}
		}
	}

	return R.list
}

//返回结果集字段
func (R *Result) Fields() map[string]string {
	return R.fields
}

//walk函数遍历
func (R *Result) Walk(callback func(v map[string]string) map[string]string ) {
	tmp := []map[string]string{}

	for _, v := range R.list {
		v = callback(v)

		if v!=nil {
			tmp = append(tmp, v) 
		}
	}

	R.list = tmp
}

//返回列表分页标准格式响应
func (R *Result) Response() map[string]any {
	return map[string]any{
		"total"      : R.Total(),
		"total_page" : R.TotalPage(),
		"page"       : R.Page,
		"num"        : R.Limit,
		"list"       : R.list,
	}
}
