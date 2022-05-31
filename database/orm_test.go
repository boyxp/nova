package database

import "testing"
import "github.com/boyxp/nova/database"

func init() {
	database.Register("database", "test", "root:123456@tcp(localhost:3306)/test")
}

func TestExec(t *testing.T) {
	sql1 := "DROP TABLE IF EXISTS goods"
	result1, err1 := database.Open("database").Exec(sql1)
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
	result2, err2 := database.Open("database").Exec(sql2)
	if err2 != nil {
		t.Log(err2)
		t.FailNow()
	}

	t.Log(result2)
}

func TestInsert(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	O.Insert(map[string]string{"name":"可口可乐","price":"100","detail":"...","category":"饮料"})
	O.Insert(map[string]string{"name":"小红帽","price":"200","detail":"...","category":"服装"})
	O.Insert(map[string]string{"name":"雪碧","price":"300","category":"饮料"})
	O.Insert(map[string]string{"name":"高跟鞋","price":"400","detail":"...","category":"服装"})
	O.Insert(map[string]string{"name":"芬达","price":"500","detail":"...","category":"饮料"})
	O.Insert(map[string]string{"name":"海魂衫","price":"600","category":"服装"})
	O.Insert(map[string]string{"name":"和其正","price":"700","detail":"...","category":"饮料"})
	O.Insert(map[string]string{"name":"领带","price":"800","detail":"...","category":"服装"})
	O.Insert(map[string]string{"name":"美年达","price":"900","category":"饮料"})
	O.Insert(map[string]string{"name":"呢子大衣","price":"1000","detail":"...","category":"服装"})
}

func TestSelectPrimary(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Where("1").Find()
	_, ok := row["name"]
	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestExist(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	ok := O.Exist("100000")
	if !ok {
		t.Log("yes")
	} else {
		t.Fail()
	}
}

func TestSelectField(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Field("left(name, 1) as c,name,category").Where("2").Find()
	_, ok := row["name"]

	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestSelectEq(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Where("goods_id", "1").Find()
	_, ok := row["name"]
	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestSelectGtEq(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Where("goods_id", ">=", "1").Find()
	_, ok := row["name"]
	if ok {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestSelectIn(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	rows := O.Where("goods_id", "in", []string{"1","2","3"}).Select()
	if len(rows)==3 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

func TestSelectNull(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Where("detail", "is", "null").Find()
	if row["detail"]=="NULL" {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestSelectNotNull(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Where("detail", "is not", "null").Find()
	if row["detail"]!="NULL" {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestSelectBetween(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	rows := O.Where("goods_id", "BETWEEN", []string{"0","3"}).Select()
	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

func TestSelectExp(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Where("goods_id>? AND detail is not null AND category=?", "3", "服装").Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestSelectLike(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Where("name", "like", "%帽%").Find()
	if len(row)>1 {
		t.Log(row)
	} else {
		t.Fail()
	}
}

func TestSelectGroup(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	rows := O.Field("count(*) as num,category,price").
			Where("goods_id", ">", "0").
			Group("category","price").
			Having("num",">",1).
			Select()

	if len(rows)>1 {
		t.Log(rows)
	} else {
		t.Fail()
	}
}

func TestSelectValue(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	name := O.Field("name").Where("goods_id", "1").Value("name")
	if len(name)>1 {
		t.Log(name)
	} else {
		t.Fail()
	}
}

func TestSelectValues(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	names := O.Field("name").Values("name")

	if len(names)>1 {
		t.Log(names)
	} else {
		t.Fail()
	}
}

func TestSelectMaxMin(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
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

func TestSelectQueryReuse(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")

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

func TestUpdate(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	af := O.Where("goods_id", "1").Update(map[string]string{"name":"可可口口","price":"111"})
	if af > 0 {
		t.Log(af)
	} else {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	af := O.Where("goods_id", "1").Delete()
	if af > 0 {
		t.Log(af)
	} else {
		t.Fail()
	}
}
