## Mysql数据库

### 1、注册数据库参数
```
//参数为数据库标识、数据库名、DSN
database.Register("database", "test", "root:123456@tcp(localhost:3306)/test")
```
### 2、创建模型

```
//参数为库名.表名，方便多个数据库的情况自主选库
Goods := database.Model{"database.goods"}
```

### 3、查询
```
Goods.Where("1").Find()
```

更多查询示例见orm_test.go

### 4、执行测试
```
go test -v
```
