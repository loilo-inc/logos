package di_test

import (
	"fmt"
	"testing"

	"github.com/loilo-inc/logos/di"
	"github.com/stretchr/testify/assert"
)

type Key string

const (
	KeyReserved1 Key = "KeyReserved1"
	KeyReserved2     = "KeyReserved2"
	KeyReserved3     = "KeyReserved3"
)

type Alice struct {
	D *di.D `di:"D"`
}

func (a *Alice) Name() string {
	return "alice"
}
func (a *Alice) Say() string {
	bob := a.D.Get("bob").(*Bob)
	return fmt.Sprintf("alice likes %s", bob.Name())
}

type Bob struct {
	D *di.D `di:"D"`
}

func (b *Bob) Name() string {
	return "bob"
}
func (b *Bob) Say() string {
	alice := b.D.Get("alice").(*Alice)
	return fmt.Sprintf("bob also likes %s", alice.Name())
}

func TestEmptyDomain(t *testing.T) {
	t.Run("should construct new object for each time", func(t *testing.T) {
		d1 := di.EmptyDomain()
		d2 := di.EmptyDomain()
		assert.NotSame(t, d1, d2)
	})
	t.Run("should panic as no value is set", func(t *testing.T) {
		assert.PanicsWithValue(t, "[di] KeyReserved1 is not registered", func() {
			d := di.EmptyDomain()
			d.Get(KeyReserved1)
		})
	})
	t.Run("can make subdomain", func(t *testing.T) {
		root := di.EmptyDomain()
		sub := root.Subdomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
		})
		v, ok := sub.Get(KeyReserved1).(int)
		assert.True(t, ok)
		assert.Equal(t, 1, v)
		assert.PanicsWithValue(t, "[di] KeyReserved1 is not registered", func() {
			root.Get(KeyReserved1)
		})
	})
}

func TestNewDomain(t *testing.T) {
	t.Run("should construct new value for each time", func(t *testing.T) {
		d1 := di.NewDomain(func(b *di.B) {})
		d2 := di.NewDomain(func(b *di.B) {})
		assert.NotSame(t, d1, d2)
	})
	t.Run("can get value", func(t *testing.T) {
		type holder struct {
			s string
		}
		h := holder{s: "holder"}
		d := di.NewDomain(func(b *di.B) {
			b.Set("int", 1)
			i := 2
			b.Set("int_ptr", &i)
			b.Set("struct", h)
			b.Set("struct_ptr", &h)
		})
		assert.Equal(t, 1, d.Get("int"))
		assert.Equal(t, 2, *(d.Get("int_ptr").(*int)))
		assert.Equal(t, h, d.Get("struct"))
		assert.Equal(t, "holder", (d.Get("struct").(holder)).s)
		assert.Equal(t, h, *(d.Get("struct_ptr").(*holder)))
	})
	t.Run("should panic on get if no value is set", func(t *testing.T) {
		assert.PanicsWithValue(t, "[di] KeyReserved1 is not registered", func() {
			d := di.NewDomain(func(b *di.B) {})
			d.Get(KeyReserved1)
		})
	})
	t.Run("can set domain as a value", func(t *testing.T) {
		type holder struct {
			D *di.D `di:"D"`
		}
		d := di.NewDomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
			b.Set(KeyReserved2, &holder{})
		})
		assert.Equal(t, 1, d.Get(KeyReserved1))
		h := d.Get(KeyReserved2).(*holder)
		assert.Same(t, d, h.D)
		assert.Equal(t, 1, h.D.Get(KeyReserved1))
	})
	t.Run("should panic if accessing bag outside of setup scope", func(t *testing.T) {
		assert.PanicsWithValue(t, "[di] domain setup func is already completed", func() {
			var bag *di.B
			di.NewDomain(func(b *di.B) {
				bag = b
			})
			bag.Set(KeyReserved1, 1)
		})
	})
	t.Run("basic", func(t *testing.T) {
		d := di.NewDomain(func(b *di.B) {
			b.Set("alice", &Alice{})
			b.Set("bob", &Bob{})
		})
		alice := d.Get("alice").(*Alice)
		bob := d.Get("bob").(*Bob)
		assert.Equal(t, d, alice.D)
		assert.Equal(t, d, bob.D)
		assert.Equal(t, "alice", alice.Name())
		assert.Equal(t, "bob", bob.Name())
		assert.Equal(t, "alice likes bob", alice.Say())
		assert.Equal(t, "bob also likes alice", bob.Say())
	})
	t.Run("should panic if value is not addressable", func(t *testing.T) {
		type holder struct {
			D *di.D `di:"D"`
		}
		assert.PanicsWithValue(t,
			"[di] panic during injection to 'holder': di_test.holder.D is not settable",
			func() {
				di.NewDomain(func(b *di.B) {
					b.Set("holder", holder{})
				})
			})
	})
	t.Run("can inject to nested pointers", func(t *testing.T) {
		type holder struct {
			D *di.D `di:"D"`
		}
		d := di.NewDomain(func(b *di.B) {
			h := &holder{}
			h2 := &h
			h3 := &h2
			b.Set("h3", h3)
		})
		h3 := d.Get("h3").(***holder)
		assert.Equal(t, d, (***h3).D)
	})
	t.Run("should panic if field for D is not public", func(t *testing.T) {
		type holder struct {
			d *di.D `di:"D"`
		}
		assert.PanicsWithValue(t,
			"[di] panic during injection to 'holder': di_test.holder.d is not settable",
			func() {
				di.NewDomain(func(b *di.B) {
					b.Set("holder", &holder{})
				})
			},
		)
	})
}

func TestD_Get(t *testing.T) {
	t.Run("can access values in parallel without mutex", func(t *testing.T) {
		d := di.NewDomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
		})
		for i := 0; i < 20; i++ {
			i := i
			t.Run(fmt.Sprintf("[%v]", i), func(t *testing.T) {
				t.Parallel()
				a, ok := d.Get(KeyReserved1).(int)
				assert.True(t, ok)
				assert.Equal(t, 1, a)
			})
		}
	})
}

func TestD_Subdomain(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		root := di.NewDomain(func(b *di.B) {})
		sub1 := root.Subdomain(func(b *di.B) {})
		assert.NotSame(t, root, sub1)
		sub2 := root.Subdomain(func(b *di.B) {})
		assert.NotSame(t, root, sub2)
		assert.NotSame(t, sub1, sub2)
	})
	t.Run("should inject Domain to value in subdomain", func(t *testing.T) {
		type holder struct {
			D *di.D `di:"D"`
		}
		root := di.NewDomain(func(b *di.B) {
			b.Set("holder_root", &holder{})
		})
		sub := root.Subdomain(func(b *di.B) {
			b.Set("holder_sub", &holder{})
		})
		assert.Equal(t, root, root.Get("holder_root").(*holder).D)
		assert.Equal(t, sub, sub.Get("holder_sub").(*holder).D)
	})
	t.Run("should panic no values registered in repo tree", func(t *testing.T) {
		root := di.NewDomain(func(b *di.B) {})
		sub := root.Subdomain(func(b *di.B) {})

		assert.PanicsWithValue(t, "[di] KeyReserved1 is not registered", func() {
			root.Get(KeyReserved1)
		})

		assert.PanicsWithValue(t, "[di] KeyReserved1 is not registered", func() {
			sub.Get(KeyReserved1)
		})
	})
	t.Run("should not affect value assignment in subdomain to parents", func(t *testing.T) {
		root := di.NewDomain(func(b *di.B) {})
		sub := root.Subdomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
		})
		v, ok := sub.Get(KeyReserved1).(int)
		assert.True(t, ok)
		assert.Equal(t, 1, v)

		assert.PanicsWithValue(t, "[di] KeyReserved1 is not registered", func() {
			root.Get(KeyReserved1)
		})
	})
	t.Run("subdomain can get values via parents", func(t *testing.T) {
		root := di.NewDomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
		})
		sub := root.Subdomain(func(b *di.B) {})

		v, ok := root.Get(KeyReserved1).(int)
		assert.True(t, ok)
		assert.Equal(t, 1, v)

		v, ok = sub.Get(KeyReserved1).(int)
		assert.True(t, ok)
		assert.Equal(t, 1, v)
	})
	t.Run("subdomain can override parent's values", func(t *testing.T) {
		root := di.NewDomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
		})
		sub := root.Subdomain(func(b *di.B) {
			b.Set(KeyReserved1, 2)
		})

		v, ok := root.Get(KeyReserved1).(int)
		assert.True(t, ok)
		assert.Equal(t, 1, v)

		v, ok = sub.Get(KeyReserved1).(int)
		assert.True(t, ok)
		assert.Equal(t, 2, v)
	})
	t.Run("can make subdomain of subdomain", func(t *testing.T) {
		parent := di.NewDomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
		})
		assert.Equal(t, 1, parent.Get(KeyReserved1))
		child := parent.Subdomain(func(b *di.B) {
			b.Set(KeyReserved2, 2)
		})
		assert.Equal(t, 1, child.Get(KeyReserved1))
		assert.Equal(t, 2, child.Get(KeyReserved2))
		grandchild := child.Subdomain(func(b *di.B) {
			b.Set(KeyReserved3, 3)
		})
		assert.Equal(t, 1, grandchild.Get(KeyReserved1))
		assert.Equal(t, 2, grandchild.Get(KeyReserved2))
		assert.Equal(t, 3, grandchild.Get(KeyReserved3))
	})
	t.Run("can set subdomain to subdomain's value", func(t *testing.T) {
		type holder struct {
			D *di.D `di:"D"`
		}
		root := di.NewDomain(func(b *di.B) {})
		sub := root.Subdomain(func(b *di.B) {
			b.Set(KeyReserved1, 1)
			b.Set(KeyReserved2, &holder{})
		})
		assert.Equal(t, 1, sub.Get(KeyReserved1))
		h := sub.Get(KeyReserved2).(*holder)
		assert.Same(t, sub, h.D)
		assert.Equal(t, 1, h.D.Get(KeyReserved1))
	})
	t.Run("should panic if accessing subdomain's bag outside of setup scope", func(t *testing.T) {
		assert.PanicsWithValue(t, "[di] domain setup func is already completed", func() {
			root := di.NewDomain(func(b *di.B) {})
			var bag *di.B
			root.Subdomain(func(b *di.B) {
				bag = b
			})
			bag.Set(KeyReserved1, 1)
		})
	})
}
