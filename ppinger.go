package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
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
	Name      string `csv:"name"`
	Hostname  string `csv:"hostname"`
	PingStats *ping.Statistics
	HttpStats bool
	hasError  bool
}

func (h *hosts) ping() {
	pinger, err := ping.NewPinger(h.Hostname)
	if err != nil {
		h.hasError = true
		return
	}

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	} else {
		pinger.SetPrivileged(false)
	}
	pinger.Count = pingCount
	pinger.Timeout = pingTimeout
	pinger.Interval = pingInterval
	pinger.Run()

	h.PingStats = pinger.Statistics()
}

func (h *hosts) httpGetRequest() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := "https://" + h.Hostname
	req, err := http.Get(url)
	if err != nil {
		return
	}
	defer req.Body.Close()

	if strings.Index(req.Status, "200") != -1 {
		h.HttpStats = true
	}
}

func (h *hosts) result() {
	if h.hasError {
		fmt.Fprintf(os.Stderr, "[ERR] %s %s\n", h.Name, h.Hostname)
	} else {
		if h.PingStats.PacketLoss == 0 || h.HttpStats == true {
			fmt.Println("[OK]", h.Hostname, h.Name, " | ", h.HttpStats)
		} else if h.PingStats.PacketLoss > 0 && h.PingStats.PacketLoss < 100 {
			fmt.Println("[WARN]", h.Hostname, h.Name, " | ", h.HttpStats)
		} else {
			fmt.Println("[ERROR]", h.Hostname, h.Name, " | ", h.HttpStats)
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
			h.httpGetRequest()
			h.result()
		}()
	}
	wg.Wait()

	fmt.Println("Press the Enter Key to terminate the console screen!")
	fmt.Scanln() // wait for Enter Key
}
