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
	return c.Context.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

//is post request
func (c *Controller) IsPost() bool {
	return c.GetRequest().Method == "POST"
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

//output json string
func (c *Controller) Json(code int64, msg string, data interface{}) {
	c.DisableView = true
	c.Context.SetHeader("Content-Type", "json")
	jsonStr, err := json.Marshal(map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
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
	return c.GetRequest().PostFormValue(field)
}

//get post array string
func (c *Controller) GetPostList(field string) []string {
	c.GetRequest().ParseForm()
	return c.GetRequest().PostForm[field]
}

//get post params
func (c *Controller) GetInputs(fields []string) map[string]string {
	values := make(map[string]string)
	c.GetRequest().ParseForm()

	for _, field := range fields {
		values[field] = c.GetInput(field)
	}
	return values
}

func (c *Controller) GetInputList(field string) []string {
	c.GetRequest().ParseForm()
	return c.GetRequest().Form[field]
}

func (c *Controller) GetInput(field string) string {
	return c.GetRequest().FormValue(field)
}

//redirct to url
func (c *Controller) Redirect(urlStr string, params map[string]string) {
	w, _ := c.GetResponse(), c.GetRequest()
	//http.Redirect(w, r, urlStr, http.StatusSeeOther)
	w.Header().Set("Location", urlStr)
	w.WriteHeader(http.StatusSeeOther)
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

//get http ResponseWriter handler
func (c *Controller) GetResponse() http.ResponseWriter {
	return c.Context.ResponseWriter
}

//get http Request handler
func (c Controller) GetRequest() *http.Request {
	return c.Context.Request
}
