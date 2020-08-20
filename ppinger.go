package main

import (
	"flag"
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
		return
	}

	pinger.Count = pingCount
	pinger.Timeout = pingTimeout
	pinger.Run()

	h.Stats = pinger.Statistics()
}

func (h *hosts) result() {
	if h.hasError {
		fmt.Println("Error:", h.Name, h.Hostname)
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

	for _, h := range hosts {
		h.ping()
		h.result()
	}
}
