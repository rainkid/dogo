package dogo

import (
	"fmt"
	"net/http"
)

var (
	Loger    = NewLoger()
	Register = NewRegister()
)

type App struct {
	Dispatcher *Dispatcher
	Host string
	Port string
}

func NewApp(host, port string) *App {
	if len(host) == 0 || len(port) == 0{
		Loger.E("start with host or port is nil")
	}
	return &App{Host:host, Port:port}
}

//app bootstrap
func (app *App) Bootstrap(router *Router) *App {
	app.Dispatcher = NewDispatcher(app, router)
	return app
}

//run application
func (app *App) Run() {
	http.Handle("/", app.Dispatcher)
	app.Listen()
}

//listen server port
func (app *App) Listen() {

	addr := fmt.Sprintf("%s:%s", app.Host, app.Port)
	Loger.I("Server started with", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		Loger.E("Server started with error : ", err.Error())
	}
}

//set app default module
func (app *App) SetDefaultModule(name string) *App {
	app.Dispatcher.SetDefaultModule(name)
	return app
}
