package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rakyll/servpprof/pprof/internal/fetch"
	"github.com/rakyll/servpprof/pprof/internal/profile"
	"github.com/rakyll/servpprof/pprof/internal/report"
	"github.com/rakyll/servpprof/pprof/internal/symbolz"
)

var (
	listen = flag.String("listen", "localhost:6464", "the hostname and port the server is listening to")
	dest   = flag.String("target", "http://localhost:6060", "the target process that enables pprof debug server")
)

var (
	// TODO(jbd): Support all profiles, including custom profiles.
	cpuRpt  = &Report{name: "profile", secs: 10}
	heapRpt = &Report{name: "heap"}
)

type Report struct {
	mu sync.Mutex
	p  *profile.Profile

	name string
	secs int
}

func (r *Report) Fetch() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// TODO(jbd): Set timeout according to the seonds parameter.
	url := fmt.Sprintf("%s/debug/pprof/%s?seconds=%s", *dest, r.name, r.secs)
	p, err := fetch.FetchProfile(url, 60*time.Second)
	if err != nil {
		return err
	}
	if err := symbolz.Symbolize(fmt.Sprintf("%s/debug/pprof/symbol", *dest), fetch.PostURL, p); err != nil {
		return err
	}
	r.p = p
	return nil
}

func (r *Report) String(cum bool) string {
	if r.p == nil {
		return ""
	}
	buf := bytes.NewBuffer(nil)
	rpt := report.NewDefault(r.p, report.Options{
		OutputFormat:   report.Text,
		CallTree:       true,
		CumSort:        cum,
		PrintAddresses: true,
	})
	report.Generate(buf, rpt, nil)
	return buf.String()
}

func (r *Report) Filter(cum bool, focus string) string {
	if r.p == nil {
		return ""
	}
	c := r.p.Copy()
	c.FilterSamplesByName(nil, nil, nil)
	buf := bytes.NewBuffer(nil)
	rpt := report.NewDefault(c, report.Options{
		OutputFormat:   report.Text,
		CallTree:       true,
		CumSort:        cum,
		PrintAddresses: true,
	})
	report.Generate(buf, rpt, nil)
	return buf.String()
}

func main() {
	cpuRpt.Fetch()
	heapRpt.Fetch()

	fmt.Println(cpuRpt.String(true))
	fmt.Println(heapRpt.String(true))
	log.Fatal(http.ListenAndServe(*listen, nil))
}
