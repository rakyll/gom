package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/rakyll/servpprof/pprof/internal/fetch"
	"github.com/rakyll/servpprof/pprof/internal/profile"
	"github.com/rakyll/servpprof/pprof/internal/report"
	"github.com/rakyll/servpprof/pprof/internal/symbolz"
	"github.com/rakyll/statik/fs"

	_ "github.com/rakyll/servpprof/statik"
)

var (
	listen = flag.String("listen", "localhost:6464", "the hostname and port the server is listening to")
	dest   = flag.String("target", "http://localhost:6060", "the target process that enables pprof debug server")
)

var (
	// TODO(jbd): Support all profiles, including custom profiles.
	reports = make(map[string]*Report)
)

type Report struct {
	mu sync.Mutex
	p  *profile.Profile

	name        string
	defaultSecs int
}

func (r *Report) Inited() bool {
	return r.p != nil
}

// Fetch fetches the current profile and the symbols from the target program.
func (r *Report) Fetch(secs int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if secs == 0 {
		secs = r.defaultSecs
	}
	// TODO(jbd): Set timeout according to the seonds parameter.
	url := fmt.Sprintf("%s/debug/pprof/%s?seconds=%d", *dest, r.name, secs)
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

// Filter filters the report with a focus regex. If no focus is provided,
// it reports back with the entire set of calls.
// Focus regex works on the package, type and function names. Filtered
// results will include parent samples from the call graph.
func (r *Report) Filter(cum bool, focus *regexp.Regexp) string {
	// TODO(jbd): Support ignore and hide regex parameters.
	if r.p == nil {
		return ""
	}
	c := r.p.Copy()
	c.FilterSamplesByName(focus, nil, nil)
	buf := bytes.NewBuffer(nil)
	rpt := report.NewDefault(c, report.Options{
		OutputFormat:   report.Text,
		CumSort:        cum,
		PrintAddresses: true,
	})
	report.Generate(buf, rpt, nil)
	// TODO(jbd): Write to a io.Writer instead.
	return buf.String()
}

func main() {
	// stats is a proxifying target/debug/pprofstats.
	// TODO(jbd): If the UI frontend knows about the target, we
	// might have eliminated the proxy handler.
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s/debug/pprofstats", *dest)
		resp, err := http.Get(url)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%v", err)
			return
		}
		defer resp.Body.Close()
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", all)
	})

	http.HandleFunc("/p", func(w http.ResponseWriter, r *http.Request) {
		p := r.FormValue("profile")
		filter := r.FormValue("filter")
		rpt, ok := reports[p]
		if !ok {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Profile not found.")
			return
		}
		if !rpt.Inited() || r.FormValue("force") == "true" {
			if err := rpt.Fetch(0); err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "%v", err)
				return
			}
		}
		if filter == "" {
			fmt.Fprint(w, rpt.Filter(true, nil))
			return
		}
		re, err := regexp.Compile(filter)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%v", err)
			return
		}
		fmt.Fprint(w, rpt.Filter(true, re))
	})

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(statikFS))

	log.Printf("Point your browser to http://%s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func init() {
	// TODO(jbd): Support user profiles.
	reports["profile"] = &Report{name: "profile", defaultSecs: 30}
	reports["heap"] = &Report{name: "heap"}
	reports["goroutine"] = &Report{name: "goroutine"}
	reports["threadcreate"] = &Report{name: "threadcreate"}
}
