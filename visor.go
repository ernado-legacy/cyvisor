package main

import (
	"flag"
	"fmt"
	"github.com/ernado/gosmsru"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	testRate    = time.Second * 5
	sendRate    = time.Minute * 10
	testTimeout = time.Second * 2
	testNumber  = "79197241488"
)

var (
	sites  = make(map[string]string)
	key    = flag.String("key", "80df3a7d-4c8c-ffb4-b197-4dc850443bba", "sms.ru key")
	group  = new(sync.WaitGroup)
	client *gosmsru.Client
)

func testUrl(url string) bool {
	res, err := http.Get(url)
	return err == nil && res.StatusCode == 200
}

func Test(url string, timeout time.Duration) bool {
	status := make(chan bool)
	go func() {
		status <- testUrl(url)
	}()

	select {
	case <-time.After(timeout):
		log.Println(url, "timed out")
		return false
	case s := <-status:
		return s
	}
}

func TestLoop(name, url string) {
	defer group.Done()
	var lastSend time.Time
	var lastStatus = true
	t := time.NewTicker(testRate)
	for _ = range t.C {
		status := Test(url, testTimeout)
		log.Println(name, status)
		if !status && time.Now().Sub(lastSend) > sendRate {
			log.Println("sending down alert")
			client.Send(testNumber, fmt.Sprintf("%s упал", name))
			lastSend = time.Now()
		}
		if status && status != lastStatus {
			log.Println("sending ok alert")
			client.Send(testNumber, fmt.Sprintf("%s поднялся", name))
			lastSend = time.Now()
		}
		lastStatus = status
	}
}

func main() {
	flag.Parse()
	fmt.Println("starting")
	client = gosmsru.New(*key)
	sites["Попутчики API"] = "http://poputchiki.cydev.ru/api/"
	sites["Попутчики nginx"] = "http://poputchiki.cydev.ru/"
	sites["Кафе.рф"] = "https://xn--80akn5b.xn--p1ai/"
	for name, url := range sites {
		group.Add(1)
		go TestLoop(name, url)
	}
	group.Wait()
}
