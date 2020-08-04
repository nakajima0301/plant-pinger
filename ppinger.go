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

type hosts struct {
	Name     string `csv:"name"`
	Hostname string `csv:"hostname"`
	Stats    *ping.Statistics
	hasError bool
}

func (h *hosts) ping() {
	pinger, err := ping.NewPinger(h.Hostname)
	if err != nil {
		fmt.Println(h.Name, "[ERROR]", h.Hostname, "no such host.")
		h.hasError = true
		return
	}

	pinger.Count = pingCount
	pinger.Timeout = pingTimeout
	pinger.Run()

	h.Stats = pinger.Statistics()
}

func (h *hosts) result() {
	if !h.hasError {
		fmt.Println(h.Name, h.Stats)
	}
}

func readDataFromCSV(filename string) []hosts {
	var h []hosts
	b, _ := ioutil.ReadFile(filename)
	csvutil.Unmarshal(b, &h)

	return h
}

func main() {
	filename := "config/default.csv"
	hosts := readDataFromCSV(filename)

	for _, h := range hosts {
		h.ping()
		h.result()
	}
}
