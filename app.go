package dogo

import (
	"fmt"
	"net/http"
	"os"
)

type App struct {
	Dispatcher *Dispatcher
	Config     *Config
	BasePath   string
}

func NewApp(file string) *App {
	config, err := NewConfig(file)
	if err != nil {
		return nil
	}
	basepath, _ := os.Getwd()
	return &App{Config: config, BasePath: basepath}
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
	port, err := app.Config.Int(app.Dispatcher.Environ, "port")

	addr := fmt.Sprintf(":%d", 8090)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("<ListenAndServe> error")
	}
	fmt.Println("listen : 127.0.0.1", port)
}

//set app default module
func (app *App) SetDefaultModule(name string) *App {
	app.Dispatcher.SetDefaultModule(name)
	return app
}
