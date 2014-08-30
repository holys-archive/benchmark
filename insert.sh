#!/bin/bash
 
LEDIS_PATH="/Users/holys/work/src/github.com/siddontang/ledisdb"
LEDIS_SERVER="ledis-server"
 
 
 
insert() {
	t1=$(date +"%s")
	echo "##########$1############\n"
	nohup ledis-server -config="$1.conf" &
    sleep 2
 
	pid=$(ps axu|grep -v grep |grep ledis-server | awk '{print $2}')
 
	for i in {1..10}
	do
		echo "$i"000w
		./benchmark -n=10000000
		ps -p $pid -o %cpu,%mem
		iostat
		echo "\n"
	done
 

	t2=$(date +"%s")
	delta=$(($t2-$t1))
	echo "$1 used $delta seconds in total\n"
	killall ledis-server
}
 
 
all() {
	insert rocksdb
#	insert  leveldb
}
 
all 
