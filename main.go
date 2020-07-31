package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/sparrc/go-ping"
)

const (
	pingCount   = 1
	pingTimeout = 3 * time.Second
)

type plant struct {
	Name     string `csv:"name"`
	Hostname string `csv:"hostname"`
}

func (p plant) ping() {
	pinger, err := ping.NewPinger(p.Hostname)
	if err != nil {
		fmt.Println(p.Name, ": [ERROR]", p.Hostname, "no such host.")
		return
	}

	pinger.Count = pingCount
	pinger.Timeout = pingTimeout
	pinger.Run()

	stats := pinger.Statistics()
	if stats.PacketLoss == 100 {
		fmt.Println(p.Name, ": [ERROR]")
	} else if stats.PacketLoss > 0 {
		fmt.Println(p.Name, ": [WARN]")
	} else {
		fmt.Println(p.Name, ": [SUCCESS]")
	}
}

func pingCsv(f string) {
	var plants []plant
	b, _ := ioutil.ReadFile(f)
	csvutil.Unmarshal(b, &plants)

	for _, p := range plants {
		p.ping()
	}
}

func main() {
	filename := "config/default.csv"
	pingCsv(filename)
}
