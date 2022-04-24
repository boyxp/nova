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
	O.Insert(map[string]interface{}{"name":"可口可乐","price":"100","detail":"...","category":"饮料"})
	O.Insert(map[string]interface{}{"name":"小红帽","price":"200","detail":"...","category":"服装"})
	O.Insert(map[string]interface{}{"name":"雪碧","price":"300","category":"饮料"})
	O.Insert(map[string]interface{}{"name":"高跟鞋","price":"400","detail":"...","category":"服装"})
	O.Insert(map[string]interface{}{"name":"芬达","price":"500","detail":"...","category":"饮料"})
	O.Insert(map[string]interface{}{"name":"海魂衫","price":"600","category":"服装"})
	O.Insert(map[string]interface{}{"name":"和其正","price":"700","detail":"...","category":"饮料"})
	O.Insert(map[string]interface{}{"name":"领带","price":"800","detail":"...","category":"服装"})
	O.Insert(map[string]interface{}{"name":"美年达","price":"900","category":"饮料"})
	O.Insert(map[string]interface{}{"name":"呢子大衣","price":"1000","detail":"...","category":"服装"})
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

func TestSelectField(t *testing.T) {
	O := (&database.Orm{}).Init("database", "goods")
	row := O.Field("price").Where("1").Find()
	_, ok := row["name"]
	if !ok {
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

//普通列表查询

//聚合列表查询
//result := o.Field("company_id,count(*) as total").Where("tax_mobile","is","null").Group("company_id").Having("total",">",1).Select()

//单字段查询
//name := o.Field("company_id  ,  receive_name,tax_mobile").Where("tax_mobile","is","null").Value("receive_name")
//fmt.Println(name)

//非聚合单项查询
//max_id := o.Field("MAX(payment_id) as max_id").Where("tax_mobile","is","null").Value("max_id")
//fmt.Println(max_id)

//min_id := o.Field("MIN(payment_id) as min_id").Where("tax_mobile","is","null").Value("min_id")
//fmt.Println(min_id)

//复用查询条件
//query := o.Field("company_id,receive_name,tax_mobile").Where("company_id","in",[]string{"1"})
//name  := query.Value("receive_name")
//fmt.Println(name)

//row   := query.Find()
//fmt.Println(row)

//count := query.Count()
//fmt.Println("count:",count)

//sum := query.Sum("payment_id")
//fmt.Println("sum:",sum)

//list := query.Select()
//for k,v := range list {
//	fmt.Println(k, v)
//}

//删除
//dar := query.Delete()
//fmt.Println("delete:", dar)

//更新
//uar := o.Where("payment_id", "11").Update(map[string]string{"tax_mobile":"138888"})
//fmt.Println("update:", uar)


