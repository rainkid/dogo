package dogo

type register struct {
	data map[string]string
}

func NewRegister() *register {
	return &register{data: make(map[string]string)}
}

func (r *register) Set(name, value string) {
	r.data[name] = value
}

func (r *register) Get(name string) string {
	return r.data[name]
}

func (r *register) Delete(name string) {
	delete(r.data, name)
}
