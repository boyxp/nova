#!/bin/bash

error(){
    echo "\033[31;5m$1\033[0m"
    exit 1
}

tip(){
    echo "\033[32;5m$1\033[0m"
}

empty=$(ls -A | wc -l)
if [ $empty != "0" ]; then
	error "请在空目录内执行初始化命令"
fi

pwd=$(pwd)
tip "当前目录为：\033[32m $pwd \033[0m \n"

if [ ! -d "/tmp/nova" ];then
	tip "1/5 首次克隆nova项目到临时目录...\n"
	git clone git@github.com:boyxp/nova.git /tmp/nova
	if [ $? -ne 0 ]; then
       error "克隆项目失败"
    fi
else
	tip "1/5 克隆nova项目到临时目录...\n"
fi

sleep 1

tip "\n2/5 正在拷贝项目结构文件...\n"

cp -r /tmp/nova/_demo/* .

tip "3/5 正在重命名模版文件...\n"
mv .env.sample .env

tip "4/5 正在安装依赖....\n"
go mod tidy

if [ $? -ne 0 ]; then
    error "安装依赖失败"
fi

tip "5/5 初始化仓库....\n"
git init
git add .
git commit -m '初始化'

tip "\n\n初始化完毕\n"
echo "运行以下命令启动项目\n"
echo "\033[35mgo run main.go\033[0m\n\n\n"
echo "然后在浏览器访问以下地址\n"
echo "127.0.0.1:9800/hello/hi?name=eve"
echo "\n\n"
