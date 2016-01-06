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
	"runtime/pprof"
	"time"

	_ "net/http/pprof"
)

type stats struct {
	Goroutine int   `json:"goroutine"`
	Thread    int   `json:"thread"`
	Block     int   `json:"block"`
	Timestamp int64 `json:"timestamp"`
}

func init() {
	http.HandleFunc("/debug/pprofstats", Stats())
}

// Stats handler returns an http.HandlerFunc that returns stats
// about the number of current goroutines, threads, etc.
// Stats handler must be accessible through "/debug/pprofstats" route
// in order for gom to display the stats from the debugged program.
func Stats() http.HandlerFunc {
	// TODO(jbd): enable block profile.
	return func(w http.ResponseWriter, r *http.Request) {
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
}
