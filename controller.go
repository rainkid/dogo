package dogo

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Controller struct {
	Context     *Context
	Dispatcher  *Dispatcher
	Layouts     []string
	DisableView bool
	Data        map[string]interface{}
	ViewPath    string
	Tpl         string
	TplExt      string
	TplFuncs    template.FuncMap
}

func (c *Controller) Construct() {
	d := c.GetDispatcher()
	c.ViewPath = "src/views"
	c.Tpl = fmt.Sprintf("%s/%s/%s", d.Module, d.Controller, strings.ToLower(d.Action))
	c.TplExt = "html"
	return
}

func (c *Controller) Destruct() {
	return
}

func (c *Controller) IsAjax() bool {
	return c.Context.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

func (c *Controller) IsPost() {

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

func (c *Controller) Json(code int64, msg string, data interface{}) {
	c.DisableView = true
	c.Context.SetHeader("Content-Type", "json")
	jsonStr, err := json.Marshal(map[string]interface{}{
		"data":    data,
		"msg":     msg,
		"success": code,
	})
	if err != nil {
		c.GetResponse().Write([]byte("json.Marshal faild."))
	}
	c.GetResponse().Write(jsonStr)
	return
}

//Render view template
func (c *Controller) Render() {
	if c.DisableView != true {
		w := c.GetResponse()
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
		err = t.Execute(w, c.Data)
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
		http.Error(c.GetResponse(), err.Error(), http.StatusInternalServerError)
	}
	return
}

func (c *Controller) GetPosts(fields []string) map[string]interface{} {
	values := make(map[string]interface{})
	c.GetRequest().ParseForm()
	for _, field := range fields {
		if value, ok := c.GetRequest().PostForm[field]; ok {
			values[field] = value
		}
	}
	return values
}

func (c *Controller) GetPost(field string) string {
	return c.GetRequest().PostFormValue(field)
}

//get post params
func (c *Controller) GetInputs(fields []string) map[string]interface{} {
	values := make(map[string]interface{})
	c.GetRequest().ParseForm()

	for _, field := range fields {
		if value, ok := c.GetRequest().Form[field]; ok {
			values[field] = value
		}
	}
	return values
}

func (c *Controller) GetInput(field string) string {
	return c.GetRequest().FormValue(field)
}

//redirct to url
func (c *Controller) Redirect(urlStr string, params map[string]string) {
	w, r := c.GetResponse(), c.GetRequest()
	http.Redirect(w, r, urlStr, http.StatusFound)
	return
}

//set cookie for http request
func (c *Controller) SetCookie(name string, value string, expire time.Duration) *Controller {
	expires := time.Now().Add(time.Second * expire)
	cookie := &http.Cookie{Name: name, Value: value, Path: "/", Expires: expires}
	http.SetCookie(c.GetResponse(), cookie)
	return c
}

//get cookie by name
func (c *Controller) GetCookie(name string) string {
	cookie, err := c.GetRequest().Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

//delete cookie by name
func (c *Controller) DelCookie(name string) *Controller {
	expires := time.Now().Add(-1)
	cookie := &http.Cookie{Name: name, Value: "value", Path: "/", Expires: expires}
	http.SetCookie(c.GetResponse(), cookie)
	return c
}

//set contorller context 
//this action will call by dispatcher
func (c *Controller) SetContext(w http.ResponseWriter, r *http.Request) *Controller {
	c.Context = NewContext(w, r)
	return c
}

//set current dispatcher
func (c *Controller) SetDispatcher(d *Dispatcher) {
	c.Dispatcher = d
}

//get current dispatcher
func (c *Controller) GetDispatcher() *Dispatcher {
	return c.Dispatcher
}

//get http ResponseWriter handler
func (c *Controller) GetResponse() http.ResponseWriter {
	return c.Context.ResponseWriter
}

//get http Request handler
func (c Controller) GetRequest() *http.Request {
	return c.Context.Request
}
