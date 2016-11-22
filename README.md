# gom [![GoDoc](https://godoc.org/github.com/rakyll/gom/http?status.svg)](https://godoc.org/github.com/rakyll/gom/http) [![Build Status](https://travis-ci.org/rakyll/gom.svg?branch=master)](https://travis-ci.org/rakyll/gom)

A visual interface to work with runtime profiling data from Go programs.

## Installation

```
go get github.com/rakyll/gom/cmd/gom
```

The program you're willing to profile should import the
github.com/rakyll/gom/http package. The http package will register several handlers to provide information about your program during runtime.

``` go
import _ "github.com/rakyll/gom/http"

// If your application is not already running an http server,
// you need to start one.
log.Println(http.ListenAndServe("localhost:6060", nil))

```

If your HTTP server is not going to handle the http.DefaultServeMux,
you need to manually register the gom handler to respond to "/debug/_gom".

For example, gorilla/mux users can use the snippet below:

``` go
import gomhttp "github.com/rakyll/gom/http"

mux := http.NewServeMux()
mux.HandleFunc("/debug/_gom", gomhttp.Handler())
log.Println(http.ListenAndServe("localhost:6060", nil))
```

Now, you are ready to launch gom.

```
$ gom
```

- :c loads the CPU profile.
- :h loads the heap profile (default profile on launch).
- :r refreshes the current profile.
- :s toggles the cumulative sort and resorts the items.
- ↓ and ↑ to paginate.
- :f=\<regex\> filters the profile with the provided regex.

## Goals

* Building a lightweight tool that works well with runtime profiles is a necessity. Over the time, I recognized that a lot of people around me delayed to use the existing pprof tools because it's a tedious experience.
* gom has no ambition to provide the features at the granularity of the features of the command line tools. Users should feel free to fallback to `go tool pprof` if they need more sophisticated features.
* Allow users to filter, hide and ignore by symbol names.
* Increase the awareness around profiling tools and packages in Go.
* Provide additional lightweight stats where possible.

### Minor Goals
* gom should provide interfaces to let the users to export their profile data and continue to work with the go tool.
* Allow users to work with their custom user profiles.
* Make it easier to generate pprof graphical output.
