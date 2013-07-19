package dogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Dispatcher struct {
	App        *App
	Router     *Router
	Module     string
	Controller string
	Action     string
	Environ    string

	Response http.ResponseWriter
	Request  *http.Request

	Found     bool
	DefModule string
}

func NewDispatcher(app *App, router *Router) *Dispatcher {
	d := &Dispatcher{App: app, Router: router}
	return d.Init()
}

//on dispatcher init 
func (d *Dispatcher) Init() *Dispatcher {
	d.Environ = "product"
	d.Module = "Index"
	d.Controller = "Index"
	d.Action = "Index"

	environ, _ := d.App.Config.String("base", "environ")
	d.Environ = environ
	return d
}

//serveHTTP
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.Found, d.Response, d.Request = false, w, r
	d.Static()
	if d.Found == false {
		d.RouteRegex()
	}
	if d.Found == false {
		d.RouteSample()
	}
	if d.Found == false {
		http.NotFound(w, r)
	}
	return
}

// route a static route
func (d *Dispatcher) Static() {
	for prefix, staticDir := range d.Router.StaticRoutes {
		if strings.HasPrefix(d.Request.URL.Path, prefix) {
			file := staticDir + d.Request.URL.Path[len(prefix):]
			http.ServeFile(d.Response, d.Request, file)
			d.Found = true
		}
	}
}

// routed an regex route
func (d *Dispatcher) RouteRegex() {
	path := d.Request.URL.Path

	if length := len(d.Router.RegexRoutes); length == 0 {
		return
	}

	for _, route := range d.Router.RegexRoutes {
		//check if Route pattern matches url
		if route.regex.MatchString(path) {
			//get submatches (params)
			matches := route.regex.FindStringSubmatch(path)

			//double check that the Route matches the URL pattern.
			if len(matches[0]) == len(path) {
				if len(route.params) > 0 {
					//add url parameters to the query param map
					values := d.Request.URL.Query()
					for i, match := range matches[1:] {
						values.Add(route.params[i], match)
					}

					//reassemble query params and add to RawQuery
					d.Request.URL.RawQuery = url.Values(values).Encode() + "&" + d.Request.URL.RawQuery
				}

				parts := strings.Split(route.forward, "/")
				d.Module, d.Controller, d.Action = parts[1], parts[2], Util_UCFirst(parts[3])
				d.Match()
				break
			}
		}
	}

}

//route a map route
func (d *Dispatcher) RouteSample() {
	path := d.Request.URL.Path

	//get request path
	parts := strings.Split(path, "/")
	length := len(parts)
	//get module, controller, action
	if length == 4 {
		d.Module = strings.ToLower(parts[1])
		d.Controller = strings.ToLower(parts[2])
		d.Action = Util_UCFirst(strings.ToLower(parts[3]))
	}
	if length < 4 && d.DefModule != "" {
		d.Module = d.DefModule

		d.Controller = strings.ToLower(parts[1])
		if length == 3 {
			d.Action = Util_UCFirst(strings.ToLower(parts[2]))
		}
		d.Request.URL.Path = strings.ToLower(fmt.Sprintf("/%s/%s/%s", d.Module, d.Controller, d.Action))
	}
	d.Match()
}

func (d *Dispatcher) Match() {
	var sampleRoute *SampleRoute
	// find maproute
	for _, value := range d.Router.SampleRoutes {
		if value.module == d.Module && value.name == d.Controller {
			sampleRoute, d.Found = value, true
			break
		}
	}

	if d.Found {
		rv := sampleRoute.NewController()
		//set context and dispatcher
		sampleRoute.CallFunc(rv, "SetContext", d.Response, d.Request)
		sampleRoute.CallFunc(rv, "SetDispatcher", d)
		//the controller contruct can not overwrite
		sampleRoute.CallFunc(rv, "Construct")

		//functions can overwrite
		sampleRoute.CallFunc(rv, "Init")
		sampleRoute.CallFunc(rv, d.Action)
		sampleRoute.CallFunc(rv, "Render")
		//on the controller destruct
		sampleRoute.CallFunc(rv, "Destruct")
	}
}

func (d *Dispatcher) SetDefaultModule(name string) *Dispatcher {
	d.DefModule = name
	return d
}
