#!/bin/bash

empty=$(ls -A | wc -l)
if [ $empty != "0" ]; then
	echo "\033[33;5m!!!请在空目录内执行初始化命令\033[0m"
	exit 1
fi

pwd=$(pwd)
echo "当前目录为\033[32m $pwd \033[0m \n"

if [ ! -d "/tmp/nova" ];then
	echo "1/5 首次克隆nova项目到临时目录...\n"
	git clone git@github.com:boyxp/nova.git /tmp/nova
else
	echo "1/5 克隆nova项目到临时目录...\n"
fi

sleep 1

echo "2/5 正在拷贝项目结构文件...\n"
cp -r /tmp/nova/_demo/ .

echo "3/5 正在重命名模版文件...\n"
mv manage.sh.sample manage.sh
mv .env.sample .env

echo "4/5 正在安装依赖....\n"
go mod tidy

echo "5/5 初始化仓库....\n"
git init
git add .
git commit -m '初始化'

echo "\033[32m\n\n初始化完毕\033[0m\n"
echo "运行以下命令启动项目\n"
echo "\033[35mgo run main.go\033[0m\n\n\n"
