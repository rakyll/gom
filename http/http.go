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

package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	httppprof "net/http/pprof"
	"runtime/pprof"
	"time"
)

type stats struct {
	Goroutine int   `json:"goroutine"`
	Thread    int   `json:"thread"`
	Block     int   `json:"block"`
	Timestamp int64 `json:"timestamp"`
}

func init() {
	AttachProfiler(http.DefaultServeMux, true)
}

// router
type router interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// AttachProfiler will register http profiling routes to http router (http.ServeMux like type that satisfy router interface)
// If you are not using http.DefaultServerMux set pprofRegistrated to false so it's http routes can also be registered.
func AttachProfiler(r router, pprofRegistrated bool) {
	r.HandleFunc("/debug/pprofstats", Stats)

	if !pprofRegistrated {
		r.HandleFunc("/debug/pprof/", httppprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", httppprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", httppprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", httppprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", httppprof.Trace)

		r.Handle("/debug/pprof/goroutine", httppprof.Handler("goroutine"))
		r.Handle("/debug/pprof/heap", httppprof.Handler("heap"))
		r.Handle("/debug/pprof/threadcreate", httppprof.Handler("threadcreate"))
		r.Handle("/debug/pprof/block", httppprof.Handler("block"))
	}
}

// Stats exposes pprof status counters, includes number of goroutines, threads, blocks
func Stats(w http.ResponseWriter, r *http.Request) {
	// TODO(jbd): enable block profile.
	n := &stats{
		Goroutine: pprof.Lookup("goroutine").Count(),
		Thread:    pprof.Lookup("threadcreate").Count(),
		Block:     pprof.Lookup("block").Count(),
		Timestamp: time.Now().Unix(),
	}
	err := json.NewEncoder(w).Encode(n)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
	}
}
