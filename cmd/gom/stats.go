package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type stats struct {
	Goroutine int   `json:"goroutine"`
	Thread    int   `json:"thread"`
	Block     int   `json:"block"`
	Timestamp int64 `json:"timestamp"`
}

func fetchStats() (s stats, err error) {
	url := fmt.Sprintf("%s/debug/pprofstats", *target)
	resp, err := http.Get(url)
	if err != nil {
		return s, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	err = d.Decode(&s)
	return
}
