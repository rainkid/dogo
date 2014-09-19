package dogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var cache = make(map[string]*SampleRoute)

type Dispatcher struct {
	App        *App
	Router     *Router
	module     string
	controller string
	action     string

	found     bool
	DefModule string
}

func NewDispatcher(app *App, router *Router) *Dispatcher {
	d := &Dispatcher{App: app, Router: router}
	return d.Init()
}

//on dispatcher init 
func (d *Dispatcher) Init() *Dispatcher {
	d.module = "Index"
	d.controller = "Index"
	d.action = "Index"
	return d
}

//serveHTTP
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.found = false
	// d.Static(w, r)
	if d.found == false {
		d.RouteRegex(w, r)
	}
	if d.found == false {
		d.RouteSample(w, r)
	}
	if d.found == false {
		http.NotFound(w, r)
	}
	return
}

// route a static route
func (d *Dispatcher) Static(w http.ResponseWriter, r *http.Request) {
	for prefix, staticDir := range d.Router.StaticRoutes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			file := staticDir + r.URL.Path[len(prefix):]
			http.ServeFile(w, r, file)
			d.found = true
		}
	}
}

// routed an regex route
func (d *Dispatcher) RouteRegex(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

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
					values := r.URL.Query()
					for i, match := range matches[1:] {
						values.Add(route.params[i], match)
					}

					//reassemble query params and add to RawQuery
					r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
				}

				parts := strings.Split(route.forward, "/")
				d.module, d.controller, d.action = parts[1], parts[2], Util_UCFirst(parts[3])
				d.Match(w, r)
				break
			}
		}
	}

}

//route a map route
func (d *Dispatcher) RouteSample(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	//get request path
	parts := strings.Split(path, "/")
	length := len(parts)
	//get module, controller, action
	if length == 4 {
		d.module = strings.ToLower(parts[1])
		d.controller = strings.ToLower(parts[2])
		d.action = Util_UCFirst(strings.ToLower(parts[3]))
	}
	if length < 4 && d.DefModule != "" {
		d.module = d.DefModule

		d.controller = strings.ToLower(parts[1])
		if length == 3 {
			d.action = Util_UCFirst(strings.ToLower(parts[2]))
		}
		r.URL.Path = strings.ToLower(fmt.Sprintf("/%s/%s/%s", d.module, d.controller, d.action))
	}
	d.Match(w, r)
}

func (d *Dispatcher) Match(w http.ResponseWriter, r *http.Request) {
	var sampleRoute *SampleRoute
	// find maproute
	for _, value := range d.Router.SampleRoutes {
		if value.module == d.module && value.name == d.controller {
			sampleRoute, d.found = value, true
			break
		}
	}

	if d.found {
		d.Exec(sampleRoute, w, r)
	}
}

func (d *Dispatcher) Exec(sampleRoute *SampleRoute, w http.ResponseWriter, r *http.Request) {
	rv := sampleRoute.NewController()
	//set context and dispatcher
	sampleRoute.CallFunc(rv, "SetContext", w, r)
	//the controller contruct can not overwrite
	sampleRoute.CallFunc(rv, "Construct", d.module, d.controller, d.action)

	//functions can overwrite
	sampleRoute.CallFunc(rv, "Init")
	sampleRoute.CallFunc(rv,  d.action)
	sampleRoute.CallFunc(rv, "Render")
	//on the controller destruct
	sampleRoute.CallFunc(rv, "Destruct")
}

func (d *Dispatcher) SetDefaultModule(name string) *Dispatcher {
	d.DefModule = name
	return d
}
