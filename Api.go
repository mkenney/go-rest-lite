/*
Package api is a Golang API service
*/
package api

import (
	"errors"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

/*
API holds the controllers for each route and the http Handler, if any
*/
type API struct {
	Controllers map[string]*Controller
	Handler     http.Handler
}

/*
NewServer returns a new API instance
*/
func NewServer() *API {
	api := new(API)
	api.Controllers = make(map[string]*Controller)
	return api
}

/*
AddHandler adds a handler to the stack
*/
func (api *API) AddHandler(endpoint string, handler func(*http.Request, *Response)) *API {
	ctrl, ok := api.Controllers[endpoint]
	if !ok {
		ctrl = NewController(endpoint)
		api.Controllers[ctrl.Endpoint] = ctrl
	}
	api.Controllers[ctrl.Endpoint].AddHandler(handler)
	return api
}

/*
GetController retrieves a controller from the stack
*/
func (api *API) GetController(endpoint string) (*Controller, error) {
	controller, ok := api.Controllers[endpoint]
	if ok {
		return controller, nil
	}
	return nil, errors.New("Controller not found")
}

/*
ListenAndServe serves all the stuff
*/
func (api *API) ListenAndServe(port int) {
	mux := http.NewServeMux()

	log.Infof("Generating endpoint handlers")
	for _, controller := range api.Controllers {
		mux.HandleFunc(controller.Endpoint, controller.HandlerFunc())
	}
	log.Infof("Done")

	log.Infof("Starting http server 0.0.0.0:%d", port)
	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}
	go log.Fatalf("%v", server.ListenAndServe())
}
