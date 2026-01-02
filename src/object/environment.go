package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment(outer *Environment) *Environment {
	s := make(map[string]Object)
	return &Environment{
		store: s,
		outer: outer,
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	if e == nil {
		return nil, false
	}

	if obj, ok := e.store[name]; ok {
		return obj, true
	}

	if e.outer != nil {
		return e.outer.Get(name)
	}

	return nil, false
}

func (e *Environment) Set(name string, value Object) (Object, bool) {
	exists := e.DoesExist(name)

	e.store[name] = value
	return value, exists
}

func (e *Environment) Declare(name string, value Object) Object {
	e.store[name] = value
	return value
}

func (e *Environment) DoesExist(name string) bool {
	_, result := e.store[name]

	return result
}

func (e *Environment) GetOuter() *Environment {
	return e.outer
}
