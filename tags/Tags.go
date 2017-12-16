package tags

import (
	"fmt"
	"strings"
)

type Tags []string

func JoinTag(group, name string) string {
	group_string := ""
	if group != "" {
		group_string = fmt.Sprintf("%s:", group)
	}
	return fmt.Sprintf("%s%s", group_string, name)
}

func SplitTag(data string) (group, name string) {
	group = ""
	name = data
	if strings.Contains(data, ":") {
		i_s := strings.SplitN(data, ":", 2)
		group = i_s[0]
		name = i_s[1]
	}
	return
}

func New(values []string) Tags {
	ret := make([]string, len(values))
	copy(ret, values)
	return ret
}

func (self Tags) HaveTag(group, name string) bool {

	string_to_check := JoinTag(group, name)

	for _, i := range self {
		if i == string_to_check {
			return true
		}
	}
	return false
}

func (self Tags) GetGroup(group string, separate bool) []string {
	ret := make([]string, 0)
	for _, i := range self {
		g, n := SplitTag(i)
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

func (self Tags) Map() map[string]([]string) {
	ret := make(map[string]([]string))
	for _, i := range self {
		g, n := SplitTag(i)

		if _, ok := ret[g]; !ok {
			ret[g] = make([]string, 0)
		}

		ret[g] = append(ret[g], n)

	}
	return ret
}

func (self Tags) Add(
	group, name string,
	abort_if_exists bool,
	remove_existing bool,
) {

	if self.HaveTag(group, name) {
		if abort_if_exists {
			return
		}

		if remove_existing {
			self.Delete(group, name)
		}
	}

	self = append(self, JoinTag(group, name))

}

func (self Tags) Delete(group, name string) {
	j := JoinTag(group, name)
	for i := len(self) - 1; i != -1; i-- {
		if self[i] == j {
			self = append(self[:i], self[i+1:]...)
		}
	}
}
