package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rakyll/servpprof/pprof"
	"github.com/rakyll/servpprof/pprof/internal/fetch"
	"github.com/rakyll/servpprof/pprof/internal/report"
	"github.com/rakyll/servpprof/pprof/internal/symbolz"
)

var target = flag.String("dest", "http://localhost:6060", "")
var listen = flag.String("listen", "localhost:6464", "the hostname and port the server is listening to")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dest := r.FormValue("dest")
		if dest == "" {
			dest = "http://localhost:6060"
		}
		url := fmt.Sprintf("%s/debug/pprof/heap?seconds=10", dest)
		p, err := fetch.FetchProfile(url, 60*time.Second)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error: %v", err)
			return
		}
		if err := symbolz.Symbolize(fmt.Sprintf("%s/debug/pprof/symbol", dest), fetch.PostURL, p); err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error: %v", err)
			return
		}

		rpt := report.NewDefault(p, report.Options{
			OutputFormat:   report.Text,
			CallTree:       true,
			CumSort:        true,
			PrintAddresses: true,
		})
		report.Generate(w, rpt, &pprof.ObjTool{})
	})

	http.HandleFunc("/sample", func(w http.ResponseWriter, r *http.Request) {
		profile := r.FormValue("profile")
		sec, _ := strconv.Atoi(r.FormValue("seconds"))
		if profile == "" {
			profile = "profile"
		}
		c := &http.Client{Timeout: time.Hour}
		url := fmt.Sprintf("%s/debug/pprof/%s?seconds=%d", *target, profile, sec)
		resp, err := c.Get(url)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error: %v", err)
			return
		}
		defer resp.Body.Close()
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error: %v", err)
			return
		}
		fmt.Fprintf(w, string(all))
	})
	log.Fatal(http.ListenAndServe(*listen, nil))
}
