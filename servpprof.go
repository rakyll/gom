package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
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
	reports = make(map[string]*Report)
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
	url := fmt.Sprintf("%s/debug/pprof/%s?seconds=%d", *dest, r.name, r.secs)
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
	http.HandleFunc("/p", func(w http.ResponseWriter, r *http.Request) {
		p := r.FormValue("profile")
		if p == "" {
			p = "heap"
		}
		rpt := reports[p]
		err := rpt.Fetch()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%v", err)
			return
		}
		fmt.Fprint(w, rpt.String(true))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.New("home").Parse(home)).Execute(w, nil)
	})
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func init() {
	reports["profile"] = &Report{name: "profile", secs: 30}
	reports["heap"] = &Report{name: "heap"}
}

var home = `<!doctype html>
<html>
<head>
  <title></title>
  <link rel="stylesheet" href="//cdn.jsdelivr.net/flat-ui/2.0/css/flat-ui.css">
  <style>
  	body {

  	}
  	.inline{display:inline}
  	.filter{width:100%}
  	.container{ width: 800px; margin: 50px auto;}
  </style>
</head>
<body>
	<div class="container">
	  <div class="row">
	  </div>
	  <div class="row">
    		<p>
    		<h3>Profiles</h3>
    		<button type="button" class="cpu btn btn-primary">CPU</button>
    		<button type="button" class="heap btn btn-primary">Heap</button>
    		</p>
    		<input type="text" class="filter" placeholder="Filter by regex...">
    		<div class="row"><pre class="results"></pre></div>
    	</div>
	</div>
	<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
	<script type="text/javascript">
		get("heap");
		$(".cpu").on("click", function() {
			get("profile");
		});
		$(".heap").on("click", function() {
			get("heap");
		});
		function get(name) {
			$('.results').html('Loading, be patient...')
			$.get('/p?profile=' + name, function(data) {
				$('.results').html(data);
			});
		};
	</script>
</body>
</html>`
