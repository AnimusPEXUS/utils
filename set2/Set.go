package set2

/*

diffs to original set:
  * support for custom comparison function
  * better safety checks on items manipulations

*/

import (
	"errors"
	"reflect"
	"sort"
)

func StringsEQFunc(i0, i1 interface{}) (bool, error) {
	for _, i := range [](interface{}){i0, i1} {
		switch i.(type) {
		case string:
		default:
			return false, errors.New("invalid type of argument")
		}
	}

	return i0.(string) == i1.(string), nil
}

type EQCheckFunc func(i0, i1 interface{}) (bool, error)

type Set struct {
	storage []interface{}
	t       reflect.Type

	fIsEQCheckFunc EQCheckFunc
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

func (self *Set) SetEQCheckFunc(f EQCheckFunc) {
	self.fIsEQCheckFunc = f
}

// func (self *Set) Storage() []interface{} {
// 	return self.storage
// }

func (self *Set) _DefineSetType(value interface{}) {
	if self.t == nil {
		self.t = reflect.TypeOf(value)
		if self.t.Kind() == reflect.String {
			self.fIsEQCheckFunc = StringsEQFunc
		}
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

func (self *Set) _Add(value interface{}) error {
	res, err := self.Have(value)

	if err != nil {
		return err
	}

	if res {
		return nil
	}

	self.storage = append(self.storage, value)

	return nil
}

func (self *Set) Add(values ...interface{}) {
	for _, val := range values {
		self._Add(val)
	}
}

func (self *Set) AddStrings(values ...string) {
	for _, val := range values {
		self.Add(val)
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

func (self *Set) Have(value interface{}) (bool, error) {
	self._DefineSetType(value)
	self._VerifySetType(value)

	for _, val := range self.storage {
		res, err := self.fIsEQCheckFunc(val, value)
		if err != nil {
			return false, err
		}

		if res {
			return true, nil
		}
	}

	return false, nil
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

func (self *Set) ListStringsSorted() []string {
	ret := self.ListStrings()
	sort.Strings(ret)
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

func (self *Set) Union(other_sets ...*Set) (*Set, error) {
	self._VerifyPairSetTypes(other_sets...)
	sets := make([]*Set, 0)
	sets = append(sets, self)
	sets = append(sets, other_sets...)
	return Union(sets...)
}

func (self *Set) Intersection(other_sets ...*Set) (*Set, error) {
	self._VerifyPairSetTypes(other_sets...)
	sets := make([]*Set, 0)
	sets = append(sets, self)
	sets = append(sets, other_sets...)
	return Intersection(sets...)
}

func (self *Set) Difference(other_sets ...*Set) (*Set, error) {
	self._VerifyPairSetTypes(other_sets...)

	new_set := new(Set)
	new_set.t = self.t

	set, err := Union(other_sets...)
	if err != nil {
		return nil, err
	}

	for _, i := range self.storage {
		res, err := set.Have(i)
		if err != nil {
			return nil, err
		}

		if !res {
			new_set.Add(i)
		}
	}

	return new_set, nil
}

func Union(sets ...*Set) (*Set, error) {

	new_set := new(Set)
	new_set.t = sets[0].t

	for _, i := range sets {
		for _, i2 := range i.storage {
			new_set.Add(i2)
		}
	}

	return new_set, nil
}

func Intersection(sets ...*Set) (*Set, error) {

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
			res, err := i2.Have(i)

			if err != nil {
				return nil, err
			}

			if !res {
				continue main_loop
			}
		}

		new_set.Add(i)
	}

	return new_set, nil
}

func Difference(sets ...*Set) (*Set, error) {
	return sets[0].Difference(sets[1:]...)
}
