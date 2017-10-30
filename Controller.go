/*
Package api provides a simple rest HTTP server
*/
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

/*
Controller stores handlers for each endpoint
*/
type Controller struct {
	Endpoint string
	Handlers []func(*http.Request, *Response)
}

/*
Done is a struct used to communicate when a goroutine is complete
*/
type Done struct{}

/*
NewController returns a new controller
*/
func NewController(endpoint string) *Controller {
	ctrl := new(Controller)
	ctrl.Endpoint = endpoint
	return ctrl
}

/*
AddHandler adds a handler to the stack
*/
func (ctrl *Controller) AddHandler(handler func(*http.Request, *Response)) *Controller {
	ctrl.Handlers = append(ctrl.Handlers, handler)
	return ctrl
}

/*
HandlerFunc returns a wrapper function that will execute all handlers in the
stack concurrently and write the results to the http response
*/
func (ctrl *Controller) HandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Print(fmt.Sprintf("Processing request %s:%s\n", request.Method, request.RequestURI))

		// Fan-out all the routines
		response := NewResponse()
		for _, handler := range ctrl.Handlers {
			go handler(request, response)
		}
		log.Print(fmt.Sprintf("Started %v routine(s)\n", len(ctrl.Handlers)))

		// Fan-in all the responses
		var responses []interface{}
		var a int
		for resp := range response.Channel {
			_, ok := resp.(handlerComplete)
			if ok {
				a++
				if a >= len(ctrl.Handlers) {
					break
				}
			} else {
				responses = append(responses, resp)
			}
		}
		log.Print(fmt.Sprintf("Collected %v response(s)\n", len(responses)))

		// Convert responses to JSON and return
		response.Body = responses
		output, err := json.Marshal(response.Body)

		if err == nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(output)
			log.Print(fmt.Sprintf("Output returned to client\n"))
		}
	}
}

func errorHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := fn(writer, request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", request.RequestURI, err)
		}
	}
}
