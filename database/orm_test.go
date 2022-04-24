package database

import "testing"
import "github.com/boyxp/nova/database"

func _TestExec(t *testing.T) {
	database.Register("database", "test", "root:123456@tcp(localhost:3306)/test")

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

//o := (&Orm{}).Init("dev", "payment")

//普通列表查询
//result := o.Where("1").Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("1").Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("company_id","1").Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("company_id",">=","1").Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("company_id","in",[]string{"1","2","3"}).Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("tax_mobile","is","null").Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("tax_mobile","is not","null").Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("company_id","BETWEEN", []string{"0","3"}).Select()
//result := o.Field("company_id,receive_name,tax_mobile").Where("company_id in (?,?) and batch_id=?", "1","2","0").Select()
//result := o.Field("payment_id,company_id,receive_name,tax_mobile").Where("receive_name","like","%尔%").Order("payment_id","desc").Page(1).Limit(5).Select()

//聚合列表查询
//result := o.Field("company_id,count(*) as total").Where("tax_mobile","is","null").Group("company_id").Having("total",">",1).Select()

//单条查询
//result := o.Field("company_id  ,  receive_name,tax_mobile").Where("tax_mobile","is","null").Find()
//for k,v := range result {
//	fmt.Println(k, v)
//}

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


