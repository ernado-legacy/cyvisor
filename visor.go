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
	key    = flag.String("key", "key", "sms.ru key")
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
		return false
	case s := <-status:
		return s
	}
}

func TestLoop(name, url string) {
	defer group.Done()
	var lastSend time.Time
	t := time.NewTicker(testRate)
	for _ = range t.C {
		status := Test(url, testTimeout)
		log.Println(name, status)
		if !status && time.Now().Sub(lastSend) > sendRate {
			log.Println("sending alert")
			err := client.Send(testNumber, fmt.Sprintf("%s упал", name))
			if err == nil {
				lastSend = time.Now()
			} else {
				log.Println(err)
			}
		}
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
