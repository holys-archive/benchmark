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

var wg sync.WaitGroup

var client *ledis.Client

var loop int = 0

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
	fmt.Printf("%s: %0.2f requests per second; consumed %f seconds\n", cmd, (float64(*number) / delta), delta)
}

func benchSet() {
	f := func() {
		waitBench("set")
	}

	bench("set", f)

}

func benchZadd() {
	f := func() {
		waitBench("zadd")
	}

	bench("zadd", f)

}

func benchSetRead() {
	println("bench set read")
}

func benchZsetRead() {
	println("bench zset read")

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
		benchZadd()
		benchSet()
	} else if *tp == "read" {
		benchSetRead()
		benchZsetRead()
	} else {
		fmt.Println("Unsupport type ")
	}

}
