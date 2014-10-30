package http

import (
	"encoding/json"
	"time"

	"net/http"
	_ "net/http/pprof"
	"runtime/pprof"
)

type stats struct {
	Goroutine int   `json:"goroutine"`
	Thread    int   `json:"thread"`
	Block     int   `json:"block"`
	Timestamp int64 `json:"timestamp"`
}

func init() {
	http.HandleFunc("/debug/pprofstats", func(w http.ResponseWriter, r *http.Request) {
		n := &stats{
			Goroutine: pprof.Lookup("goroutine").Count(),
			Thread:    pprof.Lookup("threadcreate").Count(),
			Block:     pprof.Lookup("block").Count(),
			Timestamp: time.Now().Unix(),
		}
		json.NewEncoder(w).Encode(n)
	})
}
