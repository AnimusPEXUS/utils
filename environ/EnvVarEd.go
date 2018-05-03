package environ

import (
	"fmt"
	"strings"
)

type EnvVarEd map[string]string

func New() EnvVarEd {
	return make(EnvVarEd)
}

func NewFromStrings(envs []string) EnvVarEd {
	ret := New()
	for _, i := range envs {
		res_i := strings.Index(i, "=")
		if res_i == -1 {
			panic("couldn't split string: " + i)
		}
		ret[i[:res_i]] = i[res_i+1:]
	}
	return ret
}

func (self EnvVarEd) Strings() []string {
	ret := make([]string, 0)
	for key, val := range self {
		ret = append(ret, fmt.Sprintf("%s=%s", key, val))
	}
	return ret
}

func (self EnvVarEd) Set(name string, value string) {
	self[name] = value
}

func (self EnvVarEd) Get(name string, def string) string {
	ret := def
	if t, ok := self[name]; ok {
		ret = t
	}
	return ret
}

func (self EnvVarEd) Del(name string) {
	if _, ok := self[name]; ok {
		delete(self, name)
	}
}

func (self EnvVarEd) ListNames() []string {
	ret := make([]string, 0)
	for i, _ := range self {
		ret = append(ret, i)
	}
	return ret
}

func (self EnvVarEd) UpdateWith(in EnvVarEd) {
	for k, v := range in {
		self[k] = v
	}
}
