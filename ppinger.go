package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/sparrc/go-ping"
)

const (
	pingCount    = 5
	pingTimeout  = 10 * time.Second
	pingInterval = 1 * time.Second
)

type hosts struct {
	Name     string `csv:"name"`
	Hostname string `csv:"hostname"`
	Stats    *ping.Statistics
	hasError bool
}

func (h *hosts) ping() {
	pinger, err := ping.NewPinger(h.Hostname)
	if err != nil {
		h.hasError = true
	}

	pinger.SetPrivileged(true)
	pinger.Count = pingCount
	pinger.Timeout = pingTimeout
	pinger.Interval = pingInterval
	pinger.Run()

	h.Stats = pinger.Statistics()
}

func (h *hosts) result() {

	if h.hasError {
		fmt.Println("ERROR", h.Name, h.Hostname)
	} else {
		if h.Stats.PacketLoss == 0 {
			fmt.Println("[ OK]", h.Hostname, h.Name)
		} else if h.Stats.PacketLoss > 0 && h.Stats.PacketLoss < 100 {
			fmt.Println("[WRN]", h.Hostname, h.Name)
		} else {
			fmt.Println("[ERR]", h.Hostname, h.Name)
		}
	}
}

func readDataFromCSV(filename string) []hosts {
	var h []hosts
	b, _ := ioutil.ReadFile(filename)
	csvutil.Unmarshal(b, &h)

	return h
}

func main() {
	f := flag.String("f", "config/default.csv", "Specify the configuration file")
	flag.Parse()

	hosts := readDataFromCSV(*f)

	var wg sync.WaitGroup
	for _, host := range hosts {
		wg.Add(1)
		h := host
		go func() {
			defer wg.Done()
			h.ping()
			h.result()
		}()
	}
	wg.Wait()

	fmt.Println("Press the Enter Key to terminate the console screen!")
	fmt.Scanln() // wait for Enter Key
}
