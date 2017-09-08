package dogo

import (
// "io"
)

type Output struct {
	Context *Context
	Started bool
} 

func (output *Output) Header(key, val string) {
	if output.Started == false {
		output.Context.ResponseWriter.Header().Set(key, val)
	}
}

func (output *Output) Body(content []byte) error {
		// output.Context.ResponseWriter.Started = true
	output.Started = true
	output.Context.ResponseWriter.Write(content)
	return nil
}



