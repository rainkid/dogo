package dogo

type register struct {
	data map[string]interface{}
}

func NewRegister() *register {
	return &register{data: make(map[string]interface{})}
}

func (r *register) Set(name string, i interface{}) {
	r.data[name] = i
}

func (r *register) Get(name string) interface{} {
	return r.data[name]
}

func (r *register) Delete(name string) {
	delete(r.data, name)
}
