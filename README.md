# net 

[![GoDoc](https://godoc.org/github.com/mattes/go-net?status.svg)](https://godoc.org/github.com/mattes/go-net)
[![Build Status](https://travis-ci.org/mattes/go-net.svg?branch=master)](https://travis-ci.org/mattes/go-net)


Convenience functions and wrappers around `net` and `net/http` for ...

  * Timeouts
  * Graceful Server Shutdowns with Connection Draining
  * Healthy and Ready handlers


## Example Usage 

```go
package main

import (
  "net/http"
  "github.com/mattes/go-net/server"
)

func main() {
  server.IsBehindGoogleLoadBalancer() // optional

  mux := http.NewServeMux()
  // add handlers ...
  
  // Creates a new Server with /health and /ready handlers
  // which gracefully shuts down on SIGINT or SIGTERM
  log.Fatal(server.ListenAndServe(":8080", mux))
}
```

### Health check vs Ready check

|                 | Health | Ready |
|-----------------|--------|-------|
| Running         |    200 |   200 |
| Shutting down   |    200 |   503 |


## Good Reads

* https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
* https://blog.cloudflare.com/exposing-go-on-the-internet/

