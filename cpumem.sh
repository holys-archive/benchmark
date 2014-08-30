

ps aux|grep ledis-server | head -1 |awk '{print $11, $12}'   
pid=$(ps aux|grep -v grep |grep ledis-server | awk '{print $2}')

while true
do
    date +"%Y/%m/%d %H:%M:%S"
    ps -p $pid -o %cpu,%mem
    echo "\n"
    sleep 10
done
