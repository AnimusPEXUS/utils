package tags

import (
	"fmt"
	"sort"
	"strings"
)

type Tags struct {
	values []string
}

func JoinValue(group, name string) string {
	group_string := ""
	if group != "" {
		group_string = fmt.Sprintf("%s:", group)
	}
	return fmt.Sprintf("%s%s", group_string, name)
}

func SplitValue(value string) (group, name string) {
	group = ""
	name = value
	if strings.Contains(value, ":") {
		i_s := strings.SplitN(value, ":", 2)
		group = i_s[0]
		name = i_s[1]
	}
	return
}

func New(values []string) *Tags {
	ret := new(Tags)
	ret.values = make([]string, len(values))
	copy(ret.values, values)
	return ret
}

func (self *Tags) HaveValue(value string) bool {
	for _, i := range self.values {
		if i == value {
			return true
		}
	}
	return false
}

func (self *Tags) AddValue(value string) {
	if self.HaveValue(value) {
		return
	}
	self.values = append(self.values, value)
	return
}

func (self *Tags) DeleteValue(value string) {
	for i := len(self.values) - 1; i != -1; i-- {
		if self.values[i] == value {
			self.values = append(self.values[:i], self.values[i+1:]...)
		}
	}
}

func (self *Tags) Have(group, name string) bool {
	return self.HaveValue(JoinValue(group, name))
}

func (self *Tags) HaveGroup(group string) bool {
	for _, value := range self.values {
		g, _ := SplitValue(value)
		if g == group {
			return true
		}
	}
	return false
}

func (self *Tags) Add(group, name string) {
	self.AddValue(JoinValue(group, name))
}

func (self *Tags) Delete(group, value string) {
	j := JoinValue(group, value)
	self.DeleteValue(j)
}

func (self *Tags) DeleteGroup(group string) {
	for i := len(self.values) - 1; i != -1; i-- {
		g, _ := SplitValue(self.values[i])
		if g == group {
			self.values = append(self.values[:i], self.values[i+1:]...)
		}
	}
	return
}

func (self *Tags) Values() []string {
	ret := make([]string, len(self.values))
	copy(ret, self.values)
	sort.Strings(ret)
	return ret
}

func (self *Tags) Map() map[string]([]string) {
	ret := make(map[string]([]string))
	for _, i := range self.values {
		g, n := SplitValue(i)

		if _, ok := ret[g]; !ok {
			ret[g] = make([]string, 0)
		}

		ret[g] = append(ret[g], n)

	}
	return ret
}

func (self *Tags) Get(group string, separate bool) []string {
	ret := make([]string, 0)
	for _, i := range self.values {
		g, n := SplitValue(i)
		if g == group {
			if separate {
				ret = append(ret, n)
			} else {
				ret = append(ret, i)
			}
		}
	}
	return ret
}

func (self *Tags) GetSingle(group string, separate bool) (string, bool) {
	for _, i := range self.values {
		g, n := SplitValue(i)
		if g == group {
			if separate {
				return n, true
			} else {
				return i, true
			}
		}
	}
	return "", false
}

func (self *Tags) SetSingle(group, name string) {
	self.DeleteGroup(group)
	self.Add(group, name)
	return
}

// NOTE: adding value which not contains group, will first delete all other
//       values which have no group
func (self *Tags) SetSingleValue(value string) {
	g, n := SplitValue(value)
	self.SetSingle(g, n)
	return
}
