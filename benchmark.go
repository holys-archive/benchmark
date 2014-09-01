package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/siddontang/ledisdb/client/go/ledis"
	"math/rand"
	"sync"
	"time"
)

var ip = flag.String("ip", "127.0.0.1", "ledis server ip")
var port = flag.Int("port", 6380, "ledis server port")
var clients = flag.Int("c", 50, "number of clients")
var number = flag.Int("n", 1000, "amount of keys")
var tp = flag.String("type", "write", "read or write benchmark, default write")
var count = flag.Int("count", 10000, "the amount of KV scan count")

var (
	wg     sync.WaitGroup
	client *ledis.Client
	loop   int
)

const (
    tf = "2006-01-02 15:04:05"
)

func waitBench(cmd string) {
	defer wg.Done()

	c := client.Get()
	defer c.Close()

	var err error

	v := bytes.Repeat([]byte("a"), 100)
	for i := 0; i < loop; i++ {

		n := rand.Int()

		if cmd == "zadd" {
			_, err = c.Do(cmd, "myzset", n, n)
		} else {
			_, err = c.Do(cmd, n, v)
		}

		if err != nil {
			fmt.Printf("do %s error %s", cmd, err.Error())
			return
		}
	}
}

func bench(cmd string, f func()) {
	wg.Add(*clients)

	t1 := time.Now().UnixNano()
	for i := 0; i < *clients; i++ {
		go f()
	}

	wg.Wait()

	t2 := time.Now().UnixNano()
	delta := float64(t2-t1) / float64(time.Second)
    now := fmt.Sprintf(time.Now().Format(tf))
	fmt.Printf("%s - %s: %0.2f requests per second; consumed %f seconds\n", now, cmd, (float64(*number) / delta), delta)
}

func benchSet() {
	f := func() {
		waitBench("set")
	}

	bench("set", f)

}

func benchZAdd() {
	f := func() {
		waitBench("zadd")
	}

	bench("zadd", f)

}

func waitZrangeBench(i int) {

	defer wg.Done()

	c := client.Get()
	defer c.Close()
	var start, stop int
	start = i * loop
	stop = i*loop + loop - 1
	println(start, stop)
	ay, err := ledis.Values(c.Do("zrange", "myzset", start, stop))
	if err != nil {
		fmt.Println("do zrange  err %s", err.Error())
	}

	for _, m := range ay {
		c.Do("zscore", "myzset", m.([]byte))
	}

}

func benchZSetRead() {

	c := client.Get()
	defer c.Close()

	wg.Add(*clients)

	t1 := time.Now().UnixNano()

	for i := 0; i < *clients; i++ {
		go waitZrangeBench(i)
	}

	total, _ := ledis.Int(c.Do("zcard", "myzset"))

	wg.Wait()

	t2 := time.Now().UnixNano()

	delta := float64(t2-t1) / float64(time.Second)

	if *number < total {
		total = *number
	}

	fmt.Printf("total %d zset members, consumed %f seconds\n", total, delta)

}

func setRead() {
	c := client.Get()

	var first bool = true

	var cursor []byte
	var total int64 = 0
	t1 := time.Now().UnixNano()

	for string(cursor) != "" || first {
		ay, err := ledis.Values(c.Do("scan", cursor, "count", *count))
		if err != nil {
			fmt.Printf("do scan err %s", err.Error())
			return
		} else if len(ay) != 2 {
			fmt.Println("scan result invalid")
			return
		}

		cursor = ay[0].([]byte)
		data, err := ledis.Strings(ay[1], nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, k := range data {
			ledis.String(c.Do("get", []byte(k)))
			total++

			if total%100000 == 0 {
				fmt.Println(total)
				fmt.Println(time.Now().Format("2006/01/02 15:04:05"))
			}
		}

		first = false
	}

	t2 := time.Now().UnixNano()
	delta := float64(t2-t1) / float64(time.Second)
	fmt.Printf("total %d kv keys , consumed %f seconds \n", total, delta)
}

func main() {
	flag.Parse()

	if *number <= 0 {
		panic("invalid number")
		return
	}

	if *clients <= 0 || *number < *clients {
		panic("invalid client number")
		return
	}

	loop = *number / *clients

	addr := fmt.Sprintf("%s:%d", *ip, *port)
	cfg := new(ledis.Config)
	cfg.Addr = addr
	client = ledis.NewClient(cfg)

	rand.Seed(time.Now().Unix())

	if *tp == "write" {
		benchZAdd()
		benchSet()
	} else if *tp == "read" {
		setRead()
		benchZSetRead()
	} else {
		fmt.Println("Unsupport type ")
	}

}
