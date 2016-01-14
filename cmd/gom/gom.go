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

// Package main contains a program that helps to visualize performance-related
// data from Go programs.
// Read the instructions at https://github.com/rakyll/gom to learn more.
package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	ui "github.com/gizak/termui"
)

var (
	target = flag.String("target", "http://localhost:6060", "the target process to profile; it has to enable pprof debug server")

	prompt  *ui.Par
	ls      *ui.List
	sp      *ui.Sparklines
	display *ui.Par

	cpuProfile     = &report{name: "profile"}
	heapProfile    = &report{name: "heap"}
	currentProfile = heapProfile

	promptMsg string

	reportPage  int
	reportItems []string
	cum         bool
	filter      string
)

func main() {
	flag.Parse()
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	draw()
	ui.Handle("/sys/kbd", func(e ui.Event) {
		ev := e.Data.(ui.EvtKbd)
		switch ev.KeyStr {
		case ":":
			// TODO(jbd): enable input mode and disable after esc or enter.
			promptMsg = ":"
		case "C-8":
			if l := len(promptMsg); l != 0 {
				promptMsg = promptMsg[:l-1]
			}
		case "<enter>":
			handleInput()
			promptMsg = ""
		case "<up>":
			if reportPage > 0 {
				reportPage--
			}
		case "<down>":
			reportPage++
		case "<escape>":
			promptMsg = ""
		default:
			// TODO: filter irrelevant keys such as up, down, etc.
			promptMsg += ev.KeyStr
		}
		refresh()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(ui.Event) {
		loadProfile(false)
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

func draw() {
	display = ui.NewPar("")
	display.Height = 1
	display.Border = false

	prompt = ui.NewPar(promptMsg)
	prompt.Height = 1
	prompt.Border = false

	help := ui.NewPar(`:c, :h for profiles; :f to filter; ↓ and ↑ to paginate`)
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
	sp.Height = 10
	sp.Border = false

	ls = ui.NewList()
	ls.Border = false
	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(4, 0, prompt), ui.NewCol(8, 0, help)),
		ui.NewRow(ui.NewCol(12, 0, sp)),
		ui.NewRow(ui.NewCol(12, 0, display)),
		ui.NewRow(ui.NewCol(12, 0, ls)),
	)
}

func loadStats() {
	var max = ui.TermWidth()
	s, err := fetchStats()
	if err != nil {
		displayMsg(fmt.Sprintf("error fetching stats: %v", err))
		return
	}
	displayMsg("")
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

func loadProfile(force bool) {
	if err := currentProfile.fetch(force, 0); err != nil {
		displayMsg(err.Error())
		return
	}
	re, _ := regexp.Compile(filter)
	reportItems = currentProfile.filter(cum, re)
}

func refresh() {
	prompt.Text = promptMsg

	nreport := ui.TermHeight() - 13
	ls.Height = nreport
	if len(reportItems) > nreport*reportPage {
		// can seek to the page
		ls.Items = reportItems[nreport*reportPage : len(reportItems)]
	} else {
		ls.Items = []string{}
	}

	ui.Body.Align()
	ui.Render(ui.Body)
}

func handleInput() {
	// TODO(jbd): disable input when handling input.
	switch promptMsg {
	case ":c":
		currentProfile = cpuProfile
		reportPage = 0
		filter = ""
		loadProfile(false)
	case ":h":
		currentProfile = heapProfile
		reportPage = 0
		filter = ""
		loadProfile(false)
	case ":r":
		reportPage = 0
		loadProfile(true)
	case ":s":
		cum = !cum
		reportPage = 0
		loadProfile(false)
	}
	// handle filtering
	if strings.HasPrefix(promptMsg, ":f=") {
		re := regexp.MustCompile(":f=(.*)")
		filter = re.FindStringSubmatch(promptMsg)[1]
		reportPage = 0
		loadProfile(false)
	}
	refresh()
}

func displayMsg(msg string) {
	// TODO(jbd): hide after n secs.
	display.Text = msg
	ui.Render(display)
}
