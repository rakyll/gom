// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pprof

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rakyll/servpprof/pprof/debug/gosym"

	"github.com/rakyll/servpprof/pprof/internal/fetch"
	"github.com/rakyll/servpprof/pprof/internal/plugin"
	"github.com/rakyll/servpprof/pprof/internal/profile"
	"github.com/rakyll/servpprof/pprof/internal/symbolizer"
	"github.com/rakyll/servpprof/pprof/internal/symbolz"
	"github.com/rakyll/servpprof/pprof/objfile"
)

// symbolize attempts to symbolize profile p.
// If the source is a local binary, it tries using symbolizer and obj.
// If the source is a URL, it fetches symbol information using symbolz.
func Symbolize(mode, source string, p *profile.Profile, obj plugin.ObjTool, ui plugin.UI) error {
	remote, local := true, true
	for _, o := range strings.Split(strings.ToLower(mode), ":") {
		switch o {
		case "none", "no":
			return nil
		case "local":
			remote, local = false, true
		case "remote":
			remote, local = true, false
		default:
			ui.PrintErr("ignoring unrecognized symbolization option: " + mode)
			ui.PrintErr("expecting -symbolize=[local|remote|none][:force]")
			fallthrough
		case "", "force":
			// Ignore these options, -force is recognized by symbolizer.Symbolize
		}
	}

	var err error
	if local {
		// Symbolize using binutils.
		if err = symbolizer.Symbolize(mode, p, obj, ui); err == nil {
			return nil
		}
	}
	if remote {
		err = symbolz.Symbolize(source, fetch.PostURL, p)
	}
	return err
}

// objTool implements plugin.ObjTool using Go libraries
// (instead of invoking GNU binutils).
type ObjTool struct{}

func (*ObjTool) Open(name string, start uint64) (plugin.ObjFile, error) {
	of, err := objfile.Open(name)
	if err != nil {
		return nil, err
	}
	f := &file{
		name: name,
		file: of,
	}
	return f, nil
}

func (*ObjTool) Demangle(names []string) (map[string]string, error) {
	// No C++, nothing to demangle.
	return make(map[string]string), nil
}

func (*ObjTool) Disasm(file string, start, end uint64) ([]plugin.Inst, error) {
	return nil, fmt.Errorf("disassembly not supported")
}

func (*ObjTool) SetConfig(config string) {
	// config is usually used to say what binaries to invoke.
	// Ignore entirely.
}

// file implements plugin.ObjFile using Go libraries
// (instead of invoking GNU binutils).
// A file represents a single executable being analyzed.
type file struct {
	name string
	sym  []objfile.Sym
	file *objfile.File
	pcln *gosym.Table
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Base() uint64 {
	// No support for shared libraries.
	return 0
}

func (f *file) BuildID() string {
	// No support for build ID.
	return ""
}

func (f *file) SourceLine(addr uint64) ([]plugin.Frame, error) {
	if f.pcln == nil {
		pcln, err := f.file.PCLineTable()
		if err != nil {
			return nil, err
		}
		f.pcln = pcln
	}
	file, line, fn := f.pcln.PCToLine(addr)
	if fn == nil {
		return nil, fmt.Errorf("no line information for PC=%#x", addr)
	}
	frame := []plugin.Frame{
		{
			Func: fn.Name,
			File: file,
			Line: line,
		},
	}
	return frame, nil
}

func (f *file) Symbols(r *regexp.Regexp, addr uint64) ([]*plugin.Sym, error) {
	if f.sym == nil {
		sym, err := f.file.Symbols()
		if err != nil {
			return nil, err
		}
		f.sym = sym
	}
	var out []*plugin.Sym
	for _, s := range f.sym {
		if (r == nil || r.MatchString(s.Name)) && (addr == 0 || s.Addr <= addr && addr < s.Addr+uint64(s.Size)) {
			out = append(out, &plugin.Sym{
				Name:  []string{s.Name},
				File:  f.name,
				Start: s.Addr,
				End:   s.Addr + uint64(s.Size) - 1,
			})
		}
	}
	return out, nil
}

func (f *file) Close() error {
	f.file.Close()
	return nil
}
