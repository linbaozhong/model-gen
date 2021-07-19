package tests

import (
	"fmt"
	"testing"
)

type obj struct {
	Name   string
	Depart string
	Age    int
	IS     bool
}

func (p *obj) Add(i int) {
	p.Age += i
}
func (p *obj) Del(i int) {
	p.Age -= i
}

func TestStruct(t *testing.T) {
	o := new(obj)
	o.Name = "lin"
	o.Depart = ""

	fmt.Printf("%p , %p", o, clone(o))
}

func clone(o *obj) *obj {
	a := *o
	return &a
}
