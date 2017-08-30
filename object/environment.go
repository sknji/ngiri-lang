package object

type Environment struct {
	store map[string]Object

	outer *Environment
}

func NewEnclosedEnvironment(env *Environment) *Environment {
	e := NewEnvironment()
	e.outer = env

	return e
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}
