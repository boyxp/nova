#!/bin/bash

pid() {
	if [ -f "./pid" ];then
    	PID=$(cat "./pid")
    	EXIST=$(ps aux | awk '{print $2}'| grep -w $PID)
		if [ $EXIST ];then
			echo $PID
			return
		fi
	fi

    echo 0
}

build() {
	go build api.go
	
}

if [ $# -eq 0 ];then
    echo "Usage: $0 {start|stop|status|restart}"
    exit 0
fi


case "$1" in 
	status)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "运行中...pid:$PID"
		else
			echo "未运行"
		fi
	;;
	start)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "运行中...pid:$PID"
		else
			echo "启动中..."
			build

		fi
	;;
	stop)
		PID=$(pid)
		echo $PID
#		if $(status)>0;then
#			echo "yes"
#		fi
	;;
	restart)
		echo "restart"
	;;
	*)
    	echo "Usage: $0 {start|stop|status|restart}"
        exit 1
    ;;
esac
