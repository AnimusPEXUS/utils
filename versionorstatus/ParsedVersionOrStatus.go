package versionorstatus

import (
	"fmt"
	"strconv"
	"strings"
)

// type ParsedVersion struct{ ParsedVersionOrStatus }
// type ParsedStatus struct{ ParsedVersionOrStatus }

type ParsedVersionOrStatus struct {
	values []string
	sep    string
}

func NewParsedVersionOrStatusFromString(value, sep string) *ParsedVersionOrStatus {
	self := new(ParsedVersionOrStatus)
	self.values = strings.Split(value, sep)
	self.sep = sep
	return self
}

func NewParsedVersionOrStatusFromStringSlice(
	values []string,
	sep string,
) *ParsedVersionOrStatus {
	self := new(ParsedVersionOrStatus)
	self.values = values
	self.sep = sep
	return self
}

func NewParsedVersionOrStatusFromIntSlice(value []int) *ParsedVersionOrStatus {
	self := new(ParsedVersionOrStatus)
	for _, i := range value {
		self.values = append(self.values, strconv.Itoa(i))
	}
	return self
}

func (self *ParsedVersionOrStatus) DirtyString() string {
	return strings.Join(self.values, self.sep)
}

func (self *ParsedVersionOrStatus) StrSliceString(sep string) string {
	return strings.Join(self.values, sep)
}

func (self *ParsedVersionOrStatus) IntSliceString(sep string) (string, error) {
	t, err := self.IntSlice()
	if err != nil {
		return "", err
	}

	strs := make([]string, 0)

	for _, i := range t {
		strs = append(strs, strconv.Itoa(i))
	}
	return strings.Join(strs, sep), nil
}

func (self *ParsedVersionOrStatus) InfoText() string {

	ret := fmt.Sprintf(""+
		`  Values:          %v`+
		`  StrSliceString: "%s"`,
		self.StrSlice(),
		self.StrSliceString(" "),
	)

	return ret
}

func (self *ParsedVersionOrStatus) StrSlice() []string {
	return self.values
}

func (self *ParsedVersionOrStatus) IntSlice() ([]int, error) {
	ret := make([]int, 0)
	for _, i := range self.StrSlice() {
		res, err := strconv.Atoi(i)
		if err != nil {
			return ret, err
		}
		ret = append(ret, res)
	}
	return ret, nil
}

func (self *ParsedVersionOrStatus) UIntSlice() ([]uint, error) {
	res, err := self.IntSlice()
	if err != nil {
		return nil, err
	}
	ret := make([]uint, len(res))
	for ii, i := range res {
		ret[ii] = uint(i)
	}
	return ret, nil
}
