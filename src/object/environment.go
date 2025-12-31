package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]Object
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	return val, ok
}

func (e *Environment) Set(name string, value Object) (Object, bool) {
	exists := e.DoesExist(name)

	e.store[name] = value
	return value, exists
}

func (e *Environment) Declare(name string, value Object) (Object, bool) {
	exists := e.DoesExist(name)

	e.store[name] = value
	return value, !exists
}

func (e *Environment) DoesExist(name string) bool {
	_, result := e.store[name]

	return result
}
