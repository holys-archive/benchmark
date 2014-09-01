#!/bin/bash
 
LEDIS_PATH="/Users/holys/work/src/github.com/siddontang/ledisdb"
LEDIS_SERVER="ledis-server"
 
 
 
work() {
	t1=$(date +"%s")
	echo "##########$1############\n"
	nohup ledis-server -config="$1.conf" &
    sleep 2
 
	pid=$(ps axu|grep -v grep |grep ledis-server | awk '{print $2}')
 
    echo $pid

    if [ "$2" == "write" ]; then
	    for i in {1..10}
	    do
            date +"%Y/%m/%d %H:%M:%S"
	    	echo "$i"000w
	    	./benchmark -n=10000000 -type=write

	    	ps -p $pid -o %cpu,%mem
	    	iostat
	    	echo "\n"
	    done
    elif [ "$2" == "read" ]; then
         ./benchmark -count=20000 -type=read -n=100000000
	     ps -p $pid -o %cpu,%mem
	     iostat
	     echo "\n"
    fi

	t2=$(date +"%s")
	delta=$(($t2-$t1))
	echo "$1 used $delta seconds in total\n"
	killall ledis-server
}
 
 
all() {
    work leveldb $1
    work rocksdb $1
}
 
all  $1
