package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
func (r *Report) Filter(w io.Writer, cum bool, focus *regexp.Regexp) {
	// TODO(jbd): Support ignore and hide regex parameters.
	if r.p == nil {
		return
	}
	c := r.p.Copy()
	c.FilterSamplesByName(focus, nil, nil)
	rpt := report.NewDefault(c, report.Options{
		OutputFormat:   report.JSON,
		CumSort:        cum,
		PrintAddresses: true,
	})
	report.Generate(w, rpt, nil)
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
		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(500)
			fmt.Fprint(w, all)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", all)
	})

	// p responds back with a profile report.
	http.HandleFunc("/p", func(w http.ResponseWriter, r *http.Request) {
		profile := r.FormValue("profile")
		filter := r.FormValue("filter")
		cumsort, _ := strconv.ParseBool(r.FormValue("cumsort"))
		force, _ := strconv.ParseBool(r.FormValue("force"))

		rpt, ok := reports[profile]
		if !ok {
			w.WriteHeader(404)
			fmt.Fprintf(w, "Profile not found.")
			return
		}
		if !rpt.Inited() || force {
			if err := rpt.Fetch(0); err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "%v", err)
				return
			}
		}
		var re *regexp.Regexp
		var err error
		if filter != "" {
			re, err = regexp.Compile(filter)
		}
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		rpt.Filter(w, cumsort, re)
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
