/*
Package api is a Golang API service
*/
package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

/*
Api holds the controllers for each route and the http Handler, if any
*/
type Api struct {
	Controllers map[string]*Controller
	Handler     http.Handler
}

/*
NewServer returns a new Api instance
*/
func NewServer() *Api {
	api := new(Api)
	api.Controllers = make(map[string]*Controller)
	return api
}

/*
AddHandler adds a handler to the stack
*/
func (api *Api) AddHandler(endpoint string, handler func(*http.Request, *Response)) *Api {
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
func (api *Api) GetController(endpoint string) (*Controller, error) {
	controller, ok := api.Controllers[endpoint]
	if ok {
		return controller, nil
	}
	return nil, errors.New("Controller not found")
}

/*
ListenAndServe serves all the stuff
*/
func (api *Api) ListenAndServe(port string) {
	mux := http.NewServeMux()

	fmt.Printf("generating handlers... ")
	for _, controller := range api.Controllers {
		mux.HandleFunc(controller.Endpoint, controller.HandlerFunc())
	}
	fmt.Println("done.")

	fmt.Printf("starting server on port %s\n", port)
	server := http.Server{Addr: port, Handler: mux}
	log.Fatal(server.ListenAndServe())
}
