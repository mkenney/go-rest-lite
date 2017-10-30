# Go API library experiment

Example

```golang
package main

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"time"

	"github.com/mkenney/go/api"
)

func main() {
	apiServer := api.NewServer()
	defineRoutes(apiServer)
	apiServer.ListenAndServe(":8080")
}

func defineRoutes(apiServer *api.Api) {

	// Root handler 1
	apiServer.AddHandler("/", func(request *http.Request, response *api.Response) {
		if "GET" == request.Method {
			response.Channel <- "Ok"
		}
		response.Channel <- response.Done()
	})

	// Root handler 2
	apiServer.AddHandler("/", func(request *http.Request, response *api.Response) {
		if "GET" == request.Method {
			a := make(map[string]interface{})
			a["a"] = 1
			a["b"] = "abc"
			response.Channel <- a

			rand.Seed(time.Now().UTC().UnixNano())
			bytes := make([]byte, 10)
			for i := 0; i < 10; i++ {
				bytes[i] = byte(rand.Intn(100))
			}
			response.Channel <- base64.StdEncoding.EncodeToString(bytes)
		}
		response.Channel <- response.Done()
	})

	// Root handler 3
	apiServer.AddHandler("/", func(request *http.Request, response *api.Response) {
		if "GET" == request.Method {
			a := struct {
				Thing1 string
				Thing2 int
			}{
				Thing1: "value1",
				Thing2: 2}
			response.Channel <- a
		}
		response.Channel <- response.Done()
	})

	// Foo handler 1
	apiServer.AddHandler("/foo", func(request *http.Request, response *api.Response) {
		if "GET" == request.Method {
			a := make(map[string]interface{})
			a["a"] = 3
			a["b"] = "FOO"
			response.Channel <- a
		}
		response.Channel <- response.Done()
	})
}
```
