package database

import "testing"
import "os"

//测试github工作流,docker模式,mysql8.0

//注册数据库连接
func init() {
os.Setenv("debug", "yes")
	Register("database", "test", "root:123456@tcp(localhost:3306)/test")
}

//建测试商品表
func TestExec(t *testing.T) {
	sql1 := "DROP TABLE IF EXISTS goods"
	result1, err1 := Open("database").Exec(sql1)
	if err1 != nil {
		t.Log(err1)
		t.FailNow()
	}
	t.Log(result1)

	sql2 := `CREATE TABLE goods (
	  goods_id int(11) NOT NULL AUTO_INCREMENT COMMENT '商品ID',
	  name varchar(30) NOT NULL COMMENT '商品名称',
	  price int(11) NOT NULL DEFAULT '0' COMMENT '价格单位人民币分',
	  detail varchar(100) DEFAULT NULL COMMENT '描述',
	  category varchar(100) DEFAULT NULL COMMENT '类目',
	  create_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	  update_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
	  PRIMARY KEY (goods_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8`
	result2, err2 := Open("database").Exec(sql2)
	if err2 != nil {
		t.Log(err2)
		t.FailNow()
	}

	t.Log(result2)
}

//插入
func TestInsert(t *testing.T) {
	O := Model{"goods"}
	O.Insert(map[string]string{"name":"可口可乐","price":"100","detail":"...","category":"饮料"})
	O.Insert(map[string]string{"name":"小红帽","price":"200","detail":"...","category":"服装"})
	O.Insert(map[string]string{"name":"雪碧","price":"300","category":"饮料"})
	O.Insert(map[string]string{"name":"高跟鞋","price":"400","detail":"...","category":"服装"})
	O.Insert(map[string]string{"name":"芬达","price":"500","detail":"...","category":"饮料"})
	O.Insert(map[string]string{"name":"海魂衫","price":"600","category":"服装"})
	O.Insert(map[string]string{"name":"和其正","price":"700","detail":"...","category":"饮料"})
	O.Insert(map[string]string{"name":"领带","price":"800","detail":"...","category":"服装"})
	O.Insert(map[string]string{"name":"美年达","price":"900","category":"饮料"})
	O.Insert(map[string]string{"name":"呢子大衣","price":"200","detail":"...","category":"服装"})
}

//主键条件查询
func TestSelectPrimary(t *testing.T) {
	O := Model{"goods"}
	row := O.Find("2")
	_, ok := row["name"]
	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//主键条件检查记录是否存在
func TestExist(t *testing.T) {
	O := Model{"goods"}
	ok := O.Exist("100000")
	if !ok {
		t.Log("yes")
	} else {
		t.Fail()
	}
}

//指定返回字段
func TestSelectField(t *testing.T) {
	O := Model{"goods"}
	row := O.Field("left(name, 1) as c,name,category").Where("2").Find()
	_, ok := row["name"]

	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//=条件查询
func TestSelectEq(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("goods_id", "1").Find()
	_, ok := row["name"]
	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//大于等于条件查询
func TestSelectGtEq(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("goods_id", ">=", "1").Find()
	_, ok := row["name"]
	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//多条件查询
func TestSelectMulti(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where(map[string]interface{}{
		"category":"服装",
		"name":[]interface{}{"is not", "null"},
		"price":[]interface{}{"BETWEEN", []string{"200","400"}},
	}).Select()

	if len(rows)==3 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//in条件查询
func TestSelectIn(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where("goods_id", "in", []string{"2","3","4"}).Select()
	if len(rows)==3 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//null过滤条件查询
func TestSelectNull(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail", "is", "null").Find()
	if row["detail"]=="NULL" {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//not null过滤条件查询
func TestSelectNotNull(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("detail", "is not", "null").Find()
	if row["detail"]!="NULL" {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//between区间条件查询
func TestSelectBetween(t *testing.T) {
	O := Model{"goods"}
	rows := O.Where("goods_id", "BETWEEN", []string{"0","3"}).Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//复杂语句参数代入查询
func TestSelectExp(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("goods_id>? AND detail is not null AND category=?", "3", "服装").Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//like条件搜索查询
func TestSelectLike(t *testing.T) {
	O := Model{"goods"}
	row := O.Where("name", "like", "%帽%").Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

//字段排序
func TestSelectOrder(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("category,price").Order("category","asc").Order("price","desc").Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//group+having查询
func TestSelectGroup(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("count(*) as num,category,price").
			Where("goods_id", ">", "0").
			Group("category","price").
			Having("num",">",1).
			Select()

	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//直接取单记录指定字段值
func TestSelectValue(t *testing.T) {
	O := Model{"goods"}
	name := O.Field("name").Where("goods_id", "1").Value("name")
	if len(name)>1 {
		t.Log(name)
	} else {
		t.Fail()
	}
}

//取多记录指定字段切片
func TestSelectValues(t *testing.T) {
	O := Model{"goods"}
	names := O.Field("name").Values("name")

	if len(names)>1 {
		t.Log(names)
	} else {
		t.Fail()
	}
}

//取K=>V字段记录map
func TestSelectColumns(t *testing.T) {
	O := Model{"goods"}
	//names := O.Columns("name")
	names := O.Columns("name", "goods_id")

	if len(names)>1 {
		t.Log(names)
	} else {
		t.Fail()
	}
}

//取K=>V聚合值map
func TestSelectColumnsAggs(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("count(*) as num ,category").Group("category").Columns("num", "category")

	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//取K=>row记录map
func TestSelectMap(t *testing.T) {
	O := Model{"goods"}
	rows := O.Map()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//取K=>row聚合值map
func TestSelectMapAggs(t *testing.T) {
	O := Model{"goods"}
	rows := O.Field("count(*) as num,category").Group("category").Map("num")
	if len(rows)>0 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

//取最大最小值
func TestSelectMaxMin(t *testing.T) {
	O := Model{"goods"}
	max_id := O.Field("MAX(goods_id) as max_id").Value("max_id")
	if max_id!="" {
		t.Log(max_id)
	} else {
		t.Fail()
	}

	min_id := O.Field("MIN(goods_id) as min_id").Value("min_id")
	if min_id!="" {
		t.Log(min_id)
	} else {
		t.Fail()
	}
}

//查询条件复用
func TestSelectQueryReuse(t *testing.T) {
	O := Model{"goods"}

	query := O.Where("detail", "is not", "null")

	rows := query.Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}

	row := query.Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}

	name := query.Value("name")
	if len(name)>1 {
		t.Log(name)
	} else {
		t.Fail()
	}

	count := query.Count()
	if count>1 {
		t.Log(count)
	} else {
		t.Fail()
	}

	sum := query.Sum("price")
	if sum>100 {
		t.Log(sum)
	} else {
		t.Fail()
	}
}

//更新操作，可选更新条数
func TestUpdate(t *testing.T) {
	O := Model{"goods"}
	af := O.Where("goods_id", "1").Limit(1).Update(map[string]string{"name":"可可口口","price":"111"})
	if af > 0 {
		t.Log(af)
	} else {
		t.Fail()
	}
}

//删除操作，可选删除条数
func TestDelete(t *testing.T) {
	O := Model{"goods"}
	af := O.Where("goods_id", "1").Limit(1).Delete()
	if af > 0 {
		t.Log(af)
	} else {
		t.Fail()
	}
}

//取得总条数和总页数
func TestPage(t *testing.T) {
	O := Model{"goods"}
	m := O.Limit(5)
	list  := m.Select()
	total := m.Total()
	total_page := m.TotalPage()
	if total > 0 {
		t.Log(list)
		t.Log(total)
		t.Log(total_page)
	} else {
		t.Fail()
	}
}
