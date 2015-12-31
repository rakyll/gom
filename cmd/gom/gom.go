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

	ui "github.com/gizak/termui"
)

const (
	profileCPU = iota
	profileHeap
)

var (
	target = flag.String("target", "http://localhost:6060", "the target process to profile; it has to enable pprof debug server")

	prompt *ui.Par
	sp     *ui.Sparklines

	promptMsg      string
	currentProfile *Report
	filter         string

	cpuProfile  = &Report{name: "profile", secs: 30}
	heapProfile = &Report{name: "heap"}
)

func main() {
	flag.Parse()
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	draw()
}

func draw() {
	prompt = ui.NewPar(promptMsg)
	prompt.Height = 1
	prompt.Border = false

	help := ui.NewPar(`:c, :h for profiles; :f to filter; :0 to paginate`)
	help.Height = 1
	help.Border = false
	help.TextBgColor = ui.ColorBlue
	help.Bg = ui.ColorBlue
	help.TextFgColor = ui.ColorWhite

	gs := ui.Sparkline{}
	gs.Title = "goroutines"
	gs.Height = 4
	gs.LineColor = ui.ColorCyan

	ts := ui.Sparkline{}
	ts.Title = "threads"
	ts.Height = 4
	ts.LineColor = ui.ColorCyan

	sp = ui.NewSparklines(gs, ts)
	sp.Height = 11
	sp.Border = false

	g := make([]*ui.Gauge, 10)
	for i := range g {
		g[i] = ui.NewGauge()
		g[i].LabelAlign = ui.AlignLeft
		g[i].Height = 2
		g[i].Border = false
		g[i].Percent = 100 - i*10
		g[i].PaddingBottom = 1
		g[i].BarColor = ui.ColorRed
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
	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(4, 0, prompt), ui.NewCol(8, 0, help)),
		ui.NewRow(ui.NewCol(12, 0, sp)),
		ui.NewRow(
			ui.NewCol(3, 0, g[0], g[1], g[2], g[3], g[4], g[5]),
			ui.NewCol(9, 0, ls)),
	)
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
	ui.Handle("/timer/1s", func(ui.Event) {
		loadStats()
		refresh()
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		refresh()
	})
	ui.Body.Align()
	ui.Render(ui.Body)
	ui.Loop()
}

func loadStats() {
	var max = ui.TermWidth()
	s, err := fetchStats()
	if err != nil {
		// todo: display error
		return
	}
	var cnts = []struct {
		cnt      int
		titleFmt string
	}{
		{s.Goroutine, "goroutines (%d)"},
		{s.Thread, "threads (%d)"},
	}
	for i, v := range cnts {
		if n := len(sp.Lines[i].Data); n > max {
			sp.Lines[i].Data = sp.Lines[i].Data[n-max : n]
		}
		sp.Lines[i].Title = fmt.Sprintf(v.titleFmt, v.cnt)
		sp.Lines[i].Data = append(sp.Lines[i].Data, v.cnt)
	}
}

func loadReport(force bool) {

}

func refresh() {
	prompt.Text = promptMsg
	ui.Body.Align()
	ui.Render(ui.Body)
}
