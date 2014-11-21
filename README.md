# gometry

A visual interface to work with runtime profiling data from Go programs.

![gometry screenshot](http://i.imgur.com/Wpm7VJd.png)


## Installation

```
go get github.com/rakyll/gometry/cmd/gometry
```

The program you're willing to profile should import the
github.com/rakyll/gometry/http package. The http package will register
several handlers to provide information about your program's runtime.

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

