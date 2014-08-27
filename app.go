package dogo

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	Loger    = log.New(os.Stdout, "[dogo] ", log.Ldate|log.Ltime)
	Register = NewRegister()
)

type App struct {
	Environ string
	Dispatcher *Dispatcher
	Config     *Config
	BasePath   string
}

func NewApp(file string) *App {
	config, err := NewConfig(file)
	if err != nil {
		Loger.Fatal(err.Error())
	}
	basepath, _ := os.Getwd()
	return &App{Config: config, BasePath: basepath}
}

//app bootstrap
func (app *App) Bootstrap(router *Router) *App {
	environ, _ := app.Config.String("base", "environ")
	app.Environ = environ

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
	port, err := app.Config.Int(app.Environ, "port")

	addr := fmt.Sprintf(":%d", port)
	Loger.Print("ListerAndServ ", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		Loger.Fatal("<ListenAndServe> error : ", err.Error())
	}
}

//set app default module
func (app *App) SetDefaultModule(name string) *App {
	app.Dispatcher.SetDefaultModule(name)
	return app
}
