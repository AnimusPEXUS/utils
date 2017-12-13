package set

import "reflect"

type Set struct {
	storage []interface{}
	t       reflect.Type
}

func NewSet(values ...interface{}) *Set {
	ret := new(Set)
	// ret.storage = make([]interface{}, 0)
	for _, i := range values {
		ret.Add(i)
	}
	return ret
}

func NewSetT(type_ reflect.Type) *Set {
	ret := new(Set)
	// ret.storage = make([]interface{}, 0)
	ret.t = type_
	return ret
}

func NewSetString() *Set {
	ret := new(Set)
	ret.t = reflect.TypeOf("")
	return ret
}

func NewSetInt() *Set {
	ret := new(Set)
	ret.t = reflect.TypeOf(0)
	return ret
}

// func (self *Set) Storage() []interface{} {
// 	return self.storage
// }

func (self *Set) _DefineSetType(value interface{}) {
	if self.t == nil {
		self.t = reflect.TypeOf(value)
	}
}

func (self *Set) _VerifySetType(value interface{}) {
	t := reflect.TypeOf(value)
	if t != self.t {
		panic(
			"invalid type to set. uses " + self.t.String() +
				"; requested " + t.String(),
		)
	}
}

func (self *Set) _VerifyPairSetType(other_set *Set) {
	if self.t != other_set.t {
		panic("other_set have other type")
	}
}

func (self *Set) _VerifyPairSetTypes(other_sets ...*Set) {
	for _, i := range other_sets {
		self._VerifyPairSetType(i)
	}
}

func (self *Set) _Add(value interface{}) {
	if self.Have(value) {
		return
	}

	self.storage = append(self.storage, value)
}

func (self *Set) Add(values ...interface{}) {
	for _, val := range values {
		self._Add(val)
	}
}

func (self *Set) _Remove(value interface{}) {
	self._DefineSetType(value)
	self._VerifySetType(value)

	for i, val := range self.storage {
		if val == value {
			self.storage = append(
				self.storage[:i],
				self.storage[i+1:]...,
			)
			break
		}
	}
}

func (self *Set) Remove(values ...interface{}) {
	for _, val := range values {
		self._Remove(val)
	}
}

func (self *Set) Have(value interface{}) bool {
	self._DefineSetType(value)
	self._VerifySetType(value)

	for _, val := range self.storage {
		if val == value {
			return true
		}
	}

	return false
}

func (self *Set) Len() int {
	return len(self.storage)
}

func (self *Set) List() []interface{} {
	return self.storage
}

func (self *Set) ListStrings() []string {
	ret := make([]string, 0)
	for _, val := range self.storage {
		ret = append(ret, val.(string))
	}
	return ret
}

func (self *Set) Copy() *Set {
	ret := new(Set)
	ret.t = self.t
	for _, i := range self.storage {
		ret.storage = append(ret.storage, i)
	}
	return ret
}

func (self *Set) Union(other_sets ...*Set) *Set {
	self._VerifyPairSetTypes(other_sets...)
	sets := make([]*Set, 0)
	sets = append(sets, self)
	sets = append(sets, other_sets...)
	return Union(sets...)
}

func (self *Set) Intersection(other_sets ...*Set) *Set {
	self._VerifyPairSetTypes(other_sets...)
	sets := make([]*Set, 0)
	sets = append(sets, self)
	sets = append(sets, other_sets...)
	return Intersection(sets...)
}

func (self *Set) Difference(other_sets ...*Set) *Set {
	self._VerifyPairSetTypes(other_sets...)

	new_set := new(Set)
	new_set.t = self.t

	set := Union(other_sets...)

	for _, i := range self.storage {
		if !set.Have(i) {
			new_set.Add(i)
		}
	}

	return new_set
}

func Union(sets ...*Set) *Set {

	new_set := new(Set)
	new_set.t = sets[0].t

	for _, i := range sets {
		for _, i2 := range i.storage {
			new_set.Add(i2)
		}
	}

	return new_set
}

func Intersection(sets ...*Set) *Set {

	all_values := make([]interface{}, 0)

	for _, i := range sets {
		for _, i2 := range i.storage {
			all_values = append(all_values, i2)
		}
	}

	new_set := new(Set)
	new_set.t = sets[0].t

main_loop:
	for _, i := range all_values {

		for _, i2 := range sets {
			if !i2.Have(i) {
				continue main_loop
			}
		}

		new_set.Add(i)
	}

	return new_set
}

func Difference(sets ...*Set) *Set {
	return sets[0].Difference(sets[1:]...)
}
