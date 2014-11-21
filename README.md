# gometry

A visual interface to work with runtime profiling data from Go programs.

![gometry screenshot](http://i.imgur.com/Wpm7VJd.png)


## Installation

```
go get github.com/rakyll/gometry/cmd/gometry
```

The program you're willing to profile should import the
github.com/rakyll/gometry/http package. The http package will register several handlers to provide information about your program during runtime.

``` go
import _ "github.com/rakyll/gometry/http"

// If your application is not already running an http server, you need to start one. 
go func() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
}()

```

Now, you are ready to launch the gometry.

```
$ gometry
```

Point your browser to [http://localhost:6464](http://localhost:6464).

## Goals

* Building a lightweight tool that works well with live profiles is a necessity. Over the time, I recognized that a lot of people around me delayed to use the existing pprof tools because it's a tedious experience.
* gometry has no ambition to provide the features at the granularity of the features of the command line tools. Users should feel free to fallback to `go tool pprof` if they need more sophisticated features. gometry should also provide interfaces to let the users to export their profile data and continue to work with the go tool.
* Allow users to filter, hide and ignore by symbol names. Make it easier to generate graphs, make it easier to browse extremely large image files.
* Allow users to work with their custom user profiles.
* Increase the awareness around profiling tools/packages in Go.
* Provide additional lightweight stats where possible.
