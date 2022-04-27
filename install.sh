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

go get github.com/boyxp/nova/database@v0.0.0-20220418020310-86ca24efa0b7

go get github.com/boyxp/nova@v0.0.0-20220418020310-86ca24efa0b7

echo "启动监听..."
echo "浏览器打开：127.0.0.1:9800/user/hello"

go run api.go

