## Additional steps for users of Gorilla Mux router

If you are using [Gorilla Mux router](https://github.com/gorilla/mux) you must register pprof and gom http handlers.

``` go
import (
    "net/http/pprof"
	gomhttp "github.com/rakyll/gom/http"
)

func main() {
    r := mux.NewRouter()
	attachProfiler(r)

    // ...
}

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

	router.HandleFunc("/debug/pprofstats", gomhttp.Stats)
}
```
