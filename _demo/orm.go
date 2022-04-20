package main

import (
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

func main(){
	o := Orm{table:"monitor"}
	o.Insert(map[string]interface{}{"test":"111"})

}

type Orm struct {
	db *sql.DB
	table string
}

func (O *Orm) Insert(data map[string]interface{}) {
	fmt.Println(O.table)
	O.open()

	fields       := []string{}
	placeholders := []string{}
	values       := []interface{}{}

	for k,v := range data {
		fields       = append(fields, k)
		placeholders = append(placeholders, "?")
		values       = append(values, v)
	}

	sql := "INSERT INTO "+O.table+" ("+strings.Join(fields, ",")+") VALUES("+strings.Join(placeholders, ",")+")"
	fmt.Println(sql)
	stmt, err := O.db.Prepare("INSERT INTO "+O.table+" ("+strings.Join(fields, ",")+") VALUES("+strings.Join(placeholders, ",")+")")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(values...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	O.close()
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

func (O *Orm) Field() {

}

func (O *Orm) Where() {

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

/*
//批量插入func (e *SmallormEngine) BatchInsert(data interface{}) (int64, error) {    return e.batchInsertData(data, "insert")}
//批量替换插入func (e *SmallormEngine) BatchReplace(data interface{}) (int64, error) {    return e.batchInsertData(data, "replace")}

//批量插入
func (e *SmallormEngine) batchInsertData(batchData interface{}, insertType string) (int64, error) {
  //反射解析  
  getValue := reflect.ValueOf(batchData)
  //切片大小  
  l := getValue.Len()
  //字段名  
  var fieldName []string
  //占位符  
  var placeholderString []string
  //循环判断  
  for i := 0; i < l; i++ {    value := getValue.Index(i) // Value of item    typed := value.Type()      // Type of item    if typed.Kind() != reflect.Struct {      panic("批量插入的子元素必须是结构体类型")    }
    num := value.NumField()
    //子元素值    
    var placeholder []string    //循环遍历子元素    
    for j := 0; j < num; j++ {
      //小写开头，无法反射，跳过      
      if !value.Field(j).CanInterface() {        continue      }
      //解析tag，找出真实的sql字段名      
      sqlTag := typed.Field(j).Tag.Get("sql")      
      if sqlTag != "" {        //跳过自增字段        
      	if strings.Contains(strings.ToLower(sqlTag), "auto_increment") {          continue        } 
      	else {          //字段名只记录第一个的          
      		if i == 1 {            fieldName = append(fieldName, strings.Split(sqlTag, ",")[0])          }
      		          placeholder = append(placeholder, "?")        }      } 
      		          else {        //字段名只记录第一个的        
      		          	if i == 1 {          fieldName = append(fieldName, typed.Field(j).Name)        }
      		          	        placeholder = append(placeholder, "?")      }
      //字段值      
      e.AllExec = append(e.AllExec, value.Field(j).Interface())    }
    //子元素拼接成多个()括号后的值    
    placeholderString = append(placeholderString, "("+strings.Join(placeholder, ",")+")")  }
  //拼接表，字段名，占位符  
  e.Prepare = insertType + " into " + e.GetTable() + " (" + strings.Join(fieldName, ",") + ") values " + strings.Join(placeholderString, ",")
  //prepare  
  var stmt *sql.Stmt  var err error  stmt, err = e.Db.Prepare(e.Prepare)  
  if err != nil {    return 0, e.setErrorInfo(err)  }
  //执行exec,注意这是stmt.Exec  
  result, err := stmt.Exec(e.AllExec...)  if err != nil {    return 0, e.setErrorInfo(err)  }
  //获取自增ID  
  id, _ := result.LastInsertId()  return id, nil}
*/
//自定义错误格式func (e *SmallormEngine) setErrorInfo(err error) error {  _, file, line, _ := runtime.Caller(1)  return errors.New("File: " + file + ":" + strconv.Itoa(line) + ", " + err.Error())}