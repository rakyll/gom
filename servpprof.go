package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rakyll/servpprof/pprof"
	"github.com/rakyll/servpprof/pprof/internal/fetch"
	"github.com/rakyll/servpprof/pprof/internal/report"
	"github.com/rakyll/servpprof/pprof/internal/symbolz"
)

var listen = flag.String("listen", "localhost:6464", "the hostname and port the server is listening to")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		dest := r.FormValue("dest")
		profile := r.FormValue("p")
		if dest == "" {
			dest = "http://localhost:6060"
		}
		if profile == "" {
			profile = "profile"
		}
		url := fmt.Sprintf("%s/debug/pprof/%s?seconds=%s", dest, profile, r.FormValue("seconds"))
		fmt.Println(url)
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
	log.Fatal(http.ListenAndServe(*listen, nil))
}
