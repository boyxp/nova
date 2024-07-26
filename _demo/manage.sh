#!/bin/bash

host=$(hostname)
pwd=$(pwd)
PROGRAM=$(basename $pwd)
SOURCE=main

pid() {
	if [ -f "./pid" ];then
		PID=$(cat "./pid")
		EXIST=$(ps aux | awk '{print $2}' | grep -w $PID)
		if [ $EXIST ];then
			echo $PID
			return
		fi
	fi

	echo 0
}

build() {
	go build -o $PROGRAM.$host $SOURCE.go
	if [ $? -ne 0 ];then
		echo "\033[31m build失败，启动终止 \033[0m"
		exit 1
	else
		echo "build...成功,文件名为:$PROGRAM.$host"
		return 1
	fi
}

if [ $# -eq 0 ];then
    echo "Usage: $0 {start|stop|status|restart|upgrade}"
    sh $0 status
    exit 0
fi


case "$1" in
	status|s)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "\033[32m 运行中...pid:$PID \033[0m" 
		else
			echo "\033[33m 未运行 \033[0m"
			exit 1
		fi
	;;
	start)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "\033[32m 运行中...pid:$PID \033[0m"
			exit 1
		else
			echo "启动中..."
			build
			if [ $? -ne 0 ];then
				today=$(date +%Y-%m-%d)
				nohup ./$PROGRAM.$host >> access.$host.$today.log 2>&1 &
				sleep 2
				PID=$(pid)
				if [ $PID -gt 0 ];then
					echo "\033[32m 启动成功...pid:$PID \033[0m" 
				else
					echo "\033[33m 启动失败 \033[0m"
					exit 1
				fi
			fi
		fi
	;;
	stop)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "停止中..."
			kill -s TERM $PID
			sleep 3
			PID=$(pid)
			if [ $PID -gt 0 ];then
				echo "\033[33m 停止失败...pid:$PID \033[0m"
				exit 1
			else
				echo "\033[32m 已停止 \033[0m" 
			fi
		else
			echo "\033[33m 未运行 \033[0m"
			exit 1
		fi
	;;
	restart|r)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "重启中...原pid:$PID"
			build
			if [ $? -ne 0 ];then
				kill -s HUP $PID
				sleep 2
				PID=$(pid)
				if [ $PID -gt 0 ];then
					echo "\033[32m 重启成功...新pid:$PID \033[0m" 
				else
					echo "\033[33m 重启失败 \033[0m"
					exit 1
				fi
			fi
		else
			echo "\033[33m 未运行 \033[0m"
			exit 1
		fi
	;;
	reload|l)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "重新加载配置中...原pid:$PID"
			kill -s HUP $PID
			sleep 2
			PID=$(pid)
			if [ $PID -gt 0 ];then
				echo "\033[32m 重新加载成功...新pid:$PID \033[0m" 
			else
				echo "\033[33m 重新加载失败 \033[0m"
				exit 1
			fi
		else
			echo "\033[33m 未运行 \033[0m"
			exit 1
		fi
	;;
	upgrade|u)
		output=$(git pull 2>&1)
		if [[ $output == *已经* ]]
		then
			echo "\033[33m 最近没有改动 \033[0m"
			exit 1
		fi

		if [[ $output == *up-to-date* ]]
		then
			echo "\033[33m 最近没有改动 \033[0m"
			exit 1
		fi

		echo -e "$output"

		sh $0 r
	;;
	*)
    	echo "Usage: $0 {start|stop|status|restart}"
        exit 0
    ;;
esac
