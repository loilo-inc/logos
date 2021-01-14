package di

import (
	"fmt"
)

type B struct {
	repo   map[interface{}]interface{}
	future *D
}

func (b *B) Set(key interface{}, o interface{}) {
	if b.repo == nil {
		panic("[di] domain setup func is already completed")
	}
	b.repo[key] = o
}

func (b *B) Future() *D {
	return b.future
}

type D struct {
	parent *D
	repo   map[interface{}]interface{}
}

func EmptyDomain() *D {
	return &D{
		parent: nil,
		repo:   make(map[interface{}]interface{}),
	}
}

func NewDomain(setup func(b *B)) *D {
	bag := B{
		repo:   make(map[interface{}]interface{}),
		future: new(D),
	}
	setup(&bag)
	bag.future.repo = bag.repo
	bag.repo = nil
	return bag.future
}

func (d *D) Get(key interface{}) interface{} {
	if d.repo == nil {
		panic("[di] domain is not ready yet. domain setup func are not allowed access inside domain.")
	}
	tmp := d
	for tmp != nil {
		o, ok := tmp.repo[key]
		if ok {
			return o
		}
		tmp = tmp.parent
	}
	panic(fmt.Sprintf("[di] %s is not registered", key))
}

func (d *D) Subdomain(setup func(b *B)) *D {
	if d.repo == nil {
		panic("[di] domain is not ready yet. domain setup func are not allowed access inside domain.")
	}
	bag := B{
		repo:   make(map[interface{}]interface{}),
		future: new(D),
	}
	setup(&bag)
	bag.future.parent = d
	bag.future.repo = bag.repo
	bag.repo = nil
	return bag.future
}
