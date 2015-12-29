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
	"math"

	ui "github.com/gizak/termui"
)

var (
	promptMsg = ""
	target    = flag.String("target", "http://localhost:6060", "the target process to profile; it has to enable pprof debug server")

	prompt *ui.Par
)

var reports = make(map[string]*Report)

func init() {
	reports["profile"] = &Report{name: "profile", secs: 30}
	reports["heap"] = &Report{name: "heap"}
}

func main() {
	flag.Parse()

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	draw()
}

func draw() {
	sinps := (func() []float64 {
		n := 400
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()
	data := (func() []int {
		ps := make([]int, len(sinps))
		for i, v := range sinps {
			ps[i] = int(100*v + 10)
		}
		return ps
	})()

	prompt = ui.NewPar(promptMsg)
	prompt.Height = 1
	prompt.Border = false

	help := ui.NewPar(`:c, :h for profiles; :f to filter; :0 to paginate`)
	help.Height = 1
	help.Border = false
	help.TextBgColor = ui.ColorBlue
	help.Bg = ui.ColorBlue
	help.TextFgColor = ui.ColorWhite

	goroutines := ui.Sparkline{}
	goroutines.Title = "goroutines"
	goroutines.Height = 4
	goroutines.Data = data
	goroutines.LineColor = ui.ColorCyan

	threads := ui.Sparkline{}
	threads.Title = "threads"
	threads.Height = 4
	threads.Data = data
	threads.LineColor = ui.ColorCyan

	sp := ui.NewSparklines(goroutines, threads)
	sp.Height = 11
	sp.Border = false

	gs := make([]*ui.Gauge, 10)
	for i := range gs {
		gs[i] = ui.NewGauge()
		gs[i].LabelAlign = ui.AlignLeft
		gs[i].Height = 2
		gs[i].Border = false
		gs[i].Percent = 100 - i*10
		gs[i].PaddingBottom = 1
		gs[i].BarColor = ui.ColorRed
	}

	ls := ui.NewList()
	ls.Border = false
	// ls.Items = []string{
	// 	" 0 0% 0% 0.01s 100% 00000000000105a7 runtime.notesleep",
	// 	"",
	// 	" [2] Downloading File 2",
	// 	"",
	// 	" [3] Uploading File 3",
	// 	"",
	// 	" [3] Uploading File 3",
	// 	"",
	// 	" [3] Uploading File 3",
	// }
	// ls.Height = 5

	ui.Handle("/sys/kbd", func(e ui.Event) {
		ev := e.Data.(ui.EvtKbd)
		switch ev.KeyStr {
		case ":":
			promptMsg = ":"
		case "C-8":
			if l := len(promptMsg); l != 0 {
				promptMsg = promptMsg[:l-1]
			}
		case "<enter>":
			// todo
			promptMsg = ""
		case "<escape>":
			promptMsg = ""
		default:
			promptMsg += ev.KeyStr
		}
		refresh()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(4, 0, prompt), ui.NewCol(8, 0, help)),
		ui.NewRow(ui.NewCol(12, 0, sp)),
		ui.NewRow(
			ui.NewCol(3, 0, gs[0], gs[1], gs[2], gs[3], gs[4], gs[5]),
			ui.NewCol(9, 0, ls)),
	)
	ui.Handle("1s", func(e ui.Event) {
		refresh()
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		refresh()
	})
	ui.Body.Align()
	ui.Render(ui.Body)
	ui.Loop()
}

func refresh() {
	prompt.Text = promptMsg
	ui.Body.Width = ui.TermWidth()
	ui.Body.Align()
	ui.Render(ui.Body)
}

// func statsHandler(w http.ResponseWriter, r *http.Request) {
// 	url := fmt.Sprintf("%s/debug/pprofstats", *target)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Print(err)
// 		w.WriteHeader(500)
// 		fmt.Fprintf(w, "%v", err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	all, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		w.WriteHeader(500)
// 		fmt.Fprintf(w, "%v", err)
// 		return
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		w.WriteHeader(500)
// 		fmt.Fprintf(w, "%s", all)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprintf(w, "%s", all)
// }

// func profileHandler(w http.ResponseWriter, r *http.Request) {
// 	profile := r.FormValue("profile")
// 	filter := r.FormValue("filter")
// 	img, _ := strconv.ParseBool(r.FormValue("img"))
// 	cumsort, _ := strconv.ParseBool(r.FormValue("cumsort"))
// 	force, _ := strconv.ParseBool(r.FormValue("force"))

// 	rpt, ok := reports[profile]
// 	if !ok {
// 		w.WriteHeader(404)
// 		fmt.Fprintf(w, "Profile not found.")
// 		return
// 	}
// 	if !rpt.Inited() || force {
// 		if err := rpt.Fetch(0); err != nil {
// 			w.WriteHeader(500)
// 			fmt.Fprintf(w, "%v", err)
// 			return
// 		}
// 	}
// 	var re *regexp.Regexp
// 	var err error
// 	if filter != "" {
// 		re, err = regexp.Compile(filter)
// 	}
// 	if err != nil {
// 		w.WriteHeader(400)
// 		fmt.Fprintf(w, "%v", err)
// 		return
// 	}
// 	if img {
// 		w.Header().Set("Content-Type", "image/svg+xml")
// 		rpt.Draw(w, cumsort, re)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	rpt.Filter(w, cumsort, re)
// }
