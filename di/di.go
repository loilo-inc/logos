package di

import (
	"fmt"
	"reflect"
)

type B struct {
	repo map[interface{}]interface{}
}

func (b *B) Set(key interface{}, o interface{}) {
	if b.repo == nil {
		panic("[di] domain setup func is already completed")
	}
	b.repo[key] = o
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
		repo: make(map[interface{}]interface{}),
	}
	setup(&bag)
	d := &D{}
	d.repo = bag.repo
	bag.repo = nil
	d.inject()
	return d
}

func (d *D) inject() {
	for k, v := range d.repo {
		d.injectToValue(k, v)
	}
}

func (d *D) injectToValue(k interface{}, v interface{}) {
	defer func() {
		if p := recover(); p != nil {
			panic(fmt.Sprintf("[di] panic during injection to '%s': %s", k, p))
		}
	}()
	structValue := obtainStructValue(reflect.ValueOf(v))
	if structValue == nil {
		return
	}
	t := structValue.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if _, ok := f.Tag.Lookup("di"); ok {
			field := structValue.Field(i)
			if !field.CanSet() {
				panic(fmt.Sprintf("%s.%s is not settable", t.String(), f.Name))
			} else {
				field.Set(reflect.ValueOf(d))
			}
		}
	}
}

func obtainStructValue(v reflect.Value) *reflect.Value {
	if v.Kind() == reflect.Ptr {
		return obtainStructValue(v.Elem())
	} else if v.Kind() == reflect.Struct {
		return &v
	} else {
		return nil
	}
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
	bag := &B{
		repo: make(map[interface{}]interface{}),
	}
	setup(bag)
	sd := &D{parent: d}
	sd.repo = bag.repo
	bag.repo = nil
	sd.inject()
	return sd
}
