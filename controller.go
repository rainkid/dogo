package dogo

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Controller struct {
	Context *Context

	ModuleName     string
	ControllerName string
	ActionName     string

	Layouts     []string
	DisableView bool
	Data        map[string]interface{}
	ViewPath    string
	Tpl         string
	TplExt      string
	TplFuncs    template.FuncMap
}

//@override on controller construct
func (c *Controller) BeferAction(module, controller, action string) {
	c.ModuleName, c.ControllerName, c.ActionName = module, controller, action
	c.ViewPath = "src/views"
	c.Tpl = fmt.Sprintf("%s/%s/%s", c.ModuleName, c.ControllerName, strings.ToLower(c.ActionName))
	c.TplExt = "html"
	return
}

//@override on controller destruct
func (c *Controller) AfterAction() {
	return
}

//is get request
func (c *Controller) IsAjax() bool {
	return c.Context.Input.IsAjax()
}

//is post request
func (c *Controller) IsPost() bool {
	return c.Context.Input.IsPost()
}

//Assign value to View engine
func (c *Controller) Assign(key string, value interface{}) *Controller {
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data[key] = value
	return c
}

//append layout files to layouts
func (c *Controller) Layout(file string) bool {
	tpl := fmt.Sprintf("%s/%s", c.ViewPath, file)
	if Util_FileExists(tpl) {
		c.Layouts = append(c.Layouts, tpl)
		return true
	}
	return false
}

//Render view template
func (c *Controller) Render() {
	if c.DisableView != true {
		tpl := fmt.Sprintf("%s.%s", c.Tpl, c.TplExt)

		if !c.Layout(tpl) {
			c.Error(errors.New(`can not find tpl file"` + tpl + `"`))
			return
		}
		t, err := template.New("content").Funcs(c.TplFuncs).ParseFiles(c.Layouts...)
		if err != nil {
			c.Error(err)
			return
		}
		err = t.Execute(c.Context.ResponseWriter, c.Data)
		if err != nil {
			c.Error(err)
			return
		}
	}
	return
}

//display error on reponse
func (c *Controller) Error(err error) {
	if err != nil {
		http.Error(c.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
	return
}

//get all post params that in fields
func (c *Controller) GetPosts(fields []string) map[string]string {
	values := make(map[string]string)
	for _, field := range fields {
		values[field] = c.GetPost(field)
	}
	return values
}

//get an post param
func (c *Controller) GetPost(field string) string {
	return c.Context.Request.PostFormValue(field)
}

//get post array string
func (c *Controller) GetPostList(field string) []string {
	c.Context.Request.ParseForm()
	return c.Context.Request.PostForm[field]
}

//get post params
func (c *Controller) GetInputs(fields []string) map[string]string {
	values := make(map[string]string)
	c.Context.Request.ParseForm()

	for _, field := range fields {
		values[field] = c.GetInput(field)
	}
	return values
}

func (c *Controller) GetInputList(field string) []string {
	c.Context.Request.ParseForm()
	return c.Context.Request.Form[field]
}

func (c *Controller) GetInput(field string) string {
	return c.Context.Request.FormValue(field)
}

//redirct to url
func (c *Controller) Redirect(status int, redirect string) {
	c.Context.Redirect(status, redirect)
	return
}

//set cookie for http request
func (c *Controller) SetCookie(name string, value string, expire time.Duration) *Controller {
	expires := time.Now().Add(time.Second * expire)
	cookie := &http.Cookie{Name: name, Value: value, Path: "/", Expires: expires}
	http.SetCookie(c.Context.ResponseWriter, cookie)
	return c
}

//get cookie by name
func (c *Controller) Cookie(name string) string {
	cookie, err := c.Context.Request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

//delete cookie by name
func (c *Controller) DelCookie(name string) *Controller {
	expires := time.Now().Add(-1)
	cookie := &http.Cookie{Name: name, Value: "value", Path: "/", Expires: expires}
	http.SetCookie(c.Context.ResponseWriter, cookie)
	return c
}

//set contorller context
//this action will call by dispatcher
func (c *Controller) SetContext(rw http.ResponseWriter, r *http.Request) *Controller {
	c.Context = &Context{}
	c.Context.Init(rw, r)
	return c
}
