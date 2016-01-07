// Copyright 2015 Google Inc. All Rights Reserved.
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
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rakyll/gom/internal/fetch"
	"github.com/rakyll/gom/internal/profile"
	goreport "github.com/rakyll/gom/internal/report"
	"github.com/rakyll/gom/internal/symbolz"
)

type report struct {
	mu sync.Mutex
	p  *profile.Profile

	name string
}

// fetch fetches the current profile and the symbols from the target program.
func (r *report) fetch(force bool, secs time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.p != nil && !force {
		return nil
	}
	if secs == 0 {
		secs = 60 * time.Second
	}
	url := fmt.Sprintf("%s/debug/_gom?view=profile&name=%s", *target, r.name)
	p, err := fetch.FetchProfile(url, secs)
	if err != nil {
		return err
	}
	if err := symbolz.Symbolize(fmt.Sprintf("%s/debug/_gom?view=symbol", *target), fetch.PostURL, p); err != nil {
		return err
	}
	r.p = p
	return nil
}

// filter filters the report with a focus regex. If no focus is provided,
// it reports back with the entire set of calls.
// Focus regex works on the package, type and function names. Filtered
// results will include parent samples from the call graph.
func (r *report) filter(cum bool, focus *regexp.Regexp) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.p == nil {
		return nil
	}
	c := r.p.Copy()
	c.FilterSamplesByName(focus, nil, nil)
	rpt := goreport.NewDefault(c, goreport.Options{
		OutputFormat:   goreport.Text,
		CumSort:        cum,
		PrintAddresses: true,
	})
	buf := bytes.NewBuffer(nil)
	goreport.Generate(buf, rpt, nil)
	return strings.Split(buf.String(), "\n")
}
