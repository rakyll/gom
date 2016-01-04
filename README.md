# gom

A visual interface to work with runtime profiling data from Go programs.

![gom screenshot](https://googledrive.com/host/0ByfSjdPVs9MZbkhjeUhMYzRTeEE/gom-screenshot.png)


## Installation

```
go get github.com/rakyll/gom/cmd/gom
```

The program you're willing to profile should import the
github.com/rakyll/gom/http package. The http package will register several handlers to provide information about your program during runtime.

``` go
import _ "github.com/rakyll/gom/http"

// If your application is not already running an http server, you need to start one.
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
- ← and → to paginate.
- :f=\<regex\> filters the profile with the provided regex.

### Additional installation notes

If you are using custom http.ServeMux you must register profiling http routes to your router instance.
``` go
import gomhttp "github.com/rakyll/gom/http"

mux := http.NewServeMux()
gomhttp.AttachProfiler(mux, false)
```

If you are using [Gorilla Mux router](https://github.com/gorilla/mux) when calling AttachProfiler adapter is needed.
``` go
import (
    gomhttp "github.com/rakyll/gom/http"
	"github.com/gorilla/mux"
)

r := mux.NewRouter()
gomhttp.AttachProfiler(gomhttp.GorillaMux{r}, false)
```

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
