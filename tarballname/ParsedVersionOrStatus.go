package tarballname

import (
	"fmt"
	"strconv"
	"strings"
)

type ParsedVersion struct{ ParsedVersionOrStatus }
type ParsedStatus struct{ ParsedVersionOrStatus }

type ParsedVersionOrStatus struct {
	Str      string
	DirtyStr string
	Arr      []string
	DirtyArr []string
}

func (self *ParsedVersionOrStatus) InfoText() string {
	ret := fmt.Sprintf(`  dirty_arr: %v
  dirty_str: "%s"
  arr:       %v
  str:       "%s"`,
		self.DirtyArr,
		self.DirtyStr,
		self.Arr,
		self.Str,
	)

	return ret
}

func (self *ParsedVersionOrStatus) ArrInt() ([]int, error) {
	ret := make([]int, 0)
	for _, i := range self.Arr {
		res, err := strconv.Atoi(i)
		if err != nil {
			return ret, err
		}
		ret = append(ret, res)
	}
	return ret, nil
}

func (self *ParsedVersionOrStatus) ArrUInt() ([]uint, error) {
	res, err := self.ArrInt()
	if err != nil {
		return []uint{}, err
	}
	ret := make([]uint, len(res))
	for ii, i := range res {
		ret[ii] = uint(i)
	}
	return ret, nil
}

func DefaultVersionSplitter(
	name_sliced SlicedName,
	most_possible_version Slice,
) *ParsedVersion {
	var (
		ret *ParsedVersion
	)

	ret = new(ParsedVersion)

	ret.DirtyArr = append(ret.DirtyArr[:0], ret.DirtyArr[:0]...)

	for _, j := range name_sliced[most_possible_version[0]:most_possible_version[1]] {
		ret.DirtyArr = append(ret.DirtyArr, j)
	}

	defaultVersionSplitterSub0(ret)

	return ret
}

func NewParsedVersionFromArrInt(value []int) *ParsedVersion {
	ret := new(ParsedVersion)
	for _, i := range value {
		ret.Arr = append(ret.Arr, strconv.Itoa(i))
	}

	ret.Str = strings.Join(ret.Arr, ".")

	len_arr := len(ret.Arr)
	for ii, i := range ret.Arr {
		ret.DirtyArr = append(ret.DirtyArr, i)
		if ii < len_arr-1 {
			ret.DirtyArr = append(ret.DirtyArr, ".")
		}
	}

	ret.DirtyStr = strings.Join(ret.DirtyArr, "")

	return ret
}
