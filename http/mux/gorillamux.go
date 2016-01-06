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

package mux

import (
	"net/http"

	"github.com/gorilla/mux"
)

// GorillaMux adapts GorillaMux Router structure so it can be used in place of http.ServeMux aka satisfy this package router interface
type GorillaMux struct {
	Router *mux.Router
}

// Handle passes http.Handler to underlying Gorilla Mux Router that has different function signature of Handle method.
func (r GorillaMux) Handle(path string, handler http.Handler) {
	r.Router.Handle(path, handler)
}

// HandleFunc passes http.HandlerFunc to underlying Gorilla Mux Router that has different function signature of HandleFunc method.
func (r GorillaMux) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	r.Router.HandleFunc(path, f)
}
