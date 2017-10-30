/*
Package api provides a simple rest HTTP server
*/
package api

import (
	"encoding/json"
	"log"
	"net/http"

	log "github.com/Sirupsen/logrus"
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
		log.Infof("Processing request %s:%s", request.Method, request.RequestURI)

		// Fan-out all the routines
		response := NewResponse()
		for _, handler := range ctrl.Handlers {
			go handler(request, response)
		}
		log.Infof("Started %v routine(s)", len(ctrl.Handlers))

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
		log.Infof("Collected %v response(s)", len(responses))

		// Convert responses to JSON and return
		response.Body = responses
		output, err := json.Marshal(response.Body)

		if err == nil {
			writer.Header().Set("Content-Type", "application/json")
			writer.Write(output)
			log.Infof("Output returned to client")
		}
	}
}

func errorHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := fn(writer, request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Infof("handling %q: %v", request.RequestURI, err)
		}
	}
}
