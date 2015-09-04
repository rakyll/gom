// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"

	"github.com/rakyll/statik/fs"

	_ "github.com/rakyll/gom/statik"
)

var (
	listen = flag.String("http", "localhost:6464", "the hostname and port the server is listening to")
	target = flag.String("target", "http://localhost:6060", "the target process that enables pprof debug server")
)

var reports = make(map[string]*Report)

func init() {
	reports["profile"] = &Report{name: "profile", secs: 30}
	reports["heap"] = &Report{name: "heap"}
}

func main() {
	flag.Parse()
	// stats is a proxifying target/debug/pprofstats.
	// TODO(jbd): If the UI frontend knows about the target, we
	// might have eliminated the proxy handler.
	http.HandleFunc("/stats", statsHandler)

	// p responds back with a profile report.
	http.HandleFunc("/p", profileHandler)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(statikFS))
	host, port, err := net.SplitHostPort(*listen)
	if err != nil {
		log.Fatal(err)
	}
	if host == "" {
		host = "localhost"
	}
	log.Printf("Point your browser to http://%s", net.JoinHostPort(host, port))
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/debug/pprofstats", *target)
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
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
		fmt.Fprintf(w, "%s", all)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", all)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	profile := r.FormValue("profile")
	filter := r.FormValue("filter")
	img, _ := strconv.ParseBool(r.FormValue("img"))
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
	if img {
		w.Header().Set("Content-Type", "image/svg+xml")
		rpt.Draw(w, cumsort, re)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	rpt.Filter(w, cumsort, re)
}
