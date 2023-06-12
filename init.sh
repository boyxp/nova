#!/bin/bash

if [ -d "_demo" ];then
	echo "_demo项目目录已存在"
	exit 1
fi

echo "克隆项目..."
git clone git@github.com:boyxp/nova.git


mv nova/_demo .

echo "示例文件改名..."
cd _demo

mv manage.sh.sample manage.sh

mv .env.sample .env

echo "下载go依赖..."

go mod download github.com/boyxp/nova

go get github.com/boyxp/nova/database@latest

go get github.com/boyxp/nova@latest

echo "当前项目目录为 _demo"

echo "启动监听端口：9800...(可按 Ctrl+c 终止进程)"
echo "\033[32m浏览器打开以下地址：

127.0.0.1:9800/user/hello

\033[0m"

go run main.go
