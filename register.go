package dogo 

type Data map[string]interface{}

type Register struct {
	data Data
}

func NewRegister() *Register {
	return &Register{}
}

func (r *Register) Set(name string, i interface{}) {
	r.data[name] = i;
}

func (r *Register) Get(name string) interface{}{
	return r.data[name]
}