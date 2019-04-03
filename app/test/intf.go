package test

type Intf interface {
	ChangeName()
}

type ClassA struct {
	Name string
}

type ClassB struct {
	Name string
}

func (a ClassA) ChangeName() {
	a.Name = "A has change"
}

func (a *ClassB) ChangeName() {
	a.Name = "B has change"
}

func ChangeName(h Intf) {
	h.ChangeName()
}
