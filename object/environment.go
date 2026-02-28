package object

type Environment struct {
	parent *Environment
	store  map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewSubEnvironment(parent *Environment) *Environment {
	env := NewEnvironment()
	env.parent = parent
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.parent != nil {
		obj, ok = e.parent.Get(name)
	}
	return obj, ok
}

func (e *Environment) Define(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

func (e *Environment) Assign(name string, obj Object) (Object, bool) {
	_, ok := e.Get(name)
	if !ok {
		if e.parent != nil {
			return e.parent.Assign(name, obj)
		}
		return nil, false
	}
	e.store[name] = obj
	return obj, true
}
