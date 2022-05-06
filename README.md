# Nova

基于Golang的轻量级API框架

## 特性

* 支持平滑重启
* 路由自动注册
* 编码友好的Orm查询
* 请求参数自动验证转换
* PHP开发者熟悉的异常处理
* 低侵入性设计，现有struct一行代码即可对外服务
* 一键安装快速上手
* 支持中间件

## 快速上手
创建 go.mod
```bash
module api

go 1.17
```
安装依赖
```bash
go get github.com/boyxp/nova
```
创建 hello.go
```go
package main

import "github.com/boyxp/nova"
import "github.com/boyxp/nova/router"

func main() {
   router.Register(&Hello{})
   nova.Listen("9800").Run()
}

type Hello struct {}
func (h *Hello) Hi(name string) map[string]string {
   return map[string]string{"name":"hello "+name}
}
```
启动
```bash
go run hello.go &
```

POST请求接口
```bash
curl -X POST -d 'name=eve' 127.0.0.1:9800/hello/hi
```

## 项目模式

### 初始化项目
```bash
$ curl https://raw.githubusercontent.com/boyxp/nova/master/init.sh | sh
```
默认创建 _demo 目录，可以改名为项目目录，直接初始化git

### 进程管理

启动进程
```bash
sh manage.sh start
```
查看进程状态
```bash
sh manage.sh status
```
平滑重启（重启过程旧请求不中断）
```bash
sh manage.sh restart
```

停止进程（将在请求完成后退出）
```bash
sh manage.sh stop
```

### 创建控制器
进入 controller 目录，创建struct，并将struct注册到路由

```go
package controller

import "github.com/boyxp/nova/router"
func init() {
   router.Register(&Hello{})
}

type Hello struct {}
func (h *Hello) Hi(name string) map[string]string {
	return map[string]string{"name":"hello "+name}
}
```

### 测试运行
重启进程
```bash
sh manage.sh restart
```
POST请求接口
```bash
curl -X POST -d 'name=eve' 127.0.0.1:9800/hello/hi
```
