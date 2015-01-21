package dogo

type register struct {
	data map[string]interface{}
}

func NewRegister() *register {
	return &register{data: make(map[string]interface{})}
}

func (r *register) Set(name string, value interface{}) {
	r.data[name] = value
}

func (r *register) Get(name string) interface{} {
	return r.data[name]
}

func (r *register) Delete(name string) {
	delete(r.data, name)
}
