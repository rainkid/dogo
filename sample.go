package dogo

import (
	"reflect"
	"strings"
)

type SampleRoute struct {
	module string
	name   string
	ctype  reflect.Type
}

// new sample route
func NewSampleRoute(module string, c interface{}) *SampleRoute {
	ctype := reflect.Indirect(reflect.ValueOf(c)).Type()
	name := strings.ToLower(reflect.TypeOf(c).Elem().Name())

	return &SampleRoute{module: module, name: name, ctype: ctype}
}

//call func
func (mp *SampleRoute) CallFunc(rv reflect.Value, funcName string, args ...interface{}) {
	//get SetContext method
	method := rv.MethodByName(funcName)
	//make params
	params := make([]reflect.Value, len(args))

	for i, value := range args {
		params[i] = reflect.ValueOf(value)
	}
	//call Controller method
	if method.IsValid() == true {
		method.Call(params)
	}
}

func (mp *SampleRoute) NewController() reflect.Value {
	return reflect.New(mp.ctype)
}
