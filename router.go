package dogo

type Router struct {
	SampleRoutes []*SampleRoute
	RegexRoutes  []*RegexRoute
	StaticRoutes map[string]string
}

// new a router
func NewRouter() *Router {
	return &Router{}
}

//add sample route to routers
func (router *Router) AddSampleRoute(module string, c interface{}) *Router {
	sample := NewSampleRoute(module, c)
	router.SampleRoutes = append(router.SampleRoutes, sample)
	return router
}

//add regex route to app routers
func (router *Router) AddRegexRoute(path string, forward string) *Router {
	regexRoute := NewRegexRoute(path, forward)
	router.RegexRoutes = append(router.RegexRoutes, regexRoute)
	return router
}

// add static route to app routers
func (router *Router) AddStaticRoute(dirname string, path string) *Router {
	if ok := router.StaticRoutes; ok == nil {
		router.StaticRoutes = make(map[string]string)
	}
	router.StaticRoutes[dirname] = path
	return router
}
