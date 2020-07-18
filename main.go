package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/sparrc/go-ping"
)

// Plant is ...
type Plant struct {
	Name     string `csv:"name"`
	Hostname string `csv:"hostname"`
}

// Ping is ...
func (p Plant) Ping(wg *sync.WaitGroup) {
	fmt.Println(p.Name, p.Hostname)
	pinger, err := ping.NewPinger(p.Hostname)
	if err != nil {
		panic(err)
	}

	pinger.Count = 5
	pinger.Timeout = 10 * time.Second
	pinger.Run()

	fmt.Println(pinger.Statistics())
	wg.Done()
}

// PingCsv is
func PingCsv(f string) {
	var plants []Plant
	b, _ := ioutil.ReadFile(f)
	csvutil.Unmarshal(b, &plants)

	wg := new(sync.WaitGroup)
	for _, p := range plants {
		wg.Add(1)
		go p.Ping(wg)
	}
	wg.Wait()
}

func main() {
	filename := "default.csv"
	PingCsv(filename)
}
