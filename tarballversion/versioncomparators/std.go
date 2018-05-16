package versioncomparators

import (
	"github.com/AnimusPEXUS/utils/sort"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

func init() {
	Index["std"] = &VersionComparatorStd{}
}

type VersionComparatorStd struct {
}

func (self *VersionComparatorStd) RenderNumericalVersion(
	tarballbasename *tarballname.ParsedTarballName,
) (
	[]int, error,
) {
	return tarballbasename.Version.IntSlice()
}

func (self *VersionComparatorStd) Compare(
	tarballbasename1, tarballbasename2 *tarballname.ParsedTarballName,
) (int, error) {

	one, err := tarballbasename1.Version.IntSlice()
	if err != nil {
		return -100, err
	}

	two, err := tarballbasename2.Version.IntSlice()
	if err != nil {
		return -100, err
	}

	one_c := make([]int, 0)
	for _, i := range one {
		one_c = append(one_c, i)
	}

	two_c := make([]int, 0)
	for _, i := range two {
		two_c = append(two_c, i)
	}

	for len(one_c) < len(two_c) {
		one_c = append(one_c, 0)
	}

	for len(two_c) < len(one_c) {
		two_c = append(two_c, 0)
	}

	for i := 0; i != len(one); i++ {
		if one_c[i] > two_c[i] {
			return 1, nil
		}

		if one_c[i] < two_c[i] {
			return -1, nil
		}
	}

	return 0, nil
}

func (self *VersionComparatorStd) _Sort(
	tarballbasenames_s []string,
	tarballbasenames []*tarballname.ParsedTarballName,
) error {

	what_to_sort := []interface{}{
		tarballbasenames,
	}

	if len(tarballbasenames_s) == len(tarballbasenames) {
		what_to_sort = append(what_to_sort, tarballbasenames_s)
	}

	err := sort.Sort(
		what_to_sort,
		0,
		func(
			i interface{},
			j interface{},
		) (int, error) {
			res, err := self.Compare(
				i.(*tarballname.ParsedTarballName),
				j.(*tarballname.ParsedTarballName),
			)
			if err != nil {
				return -100, err
			}
			return res, nil
		},
		false,
	)
	if err != nil {
		return err
	}

	return nil
}

func (self *VersionComparatorStd) Sort(
	tarballbasenames []*tarballname.ParsedTarballName,
) error {
	return self._Sort([]string{}, tarballbasenames)
}

func (self *VersionComparatorStd) SortStrings(
	tarballbasenames_s []string,
	parser types.TarballNameParserI,
) error {

	tarballbasenames := make([]*tarballname.ParsedTarballName, 0)
	for _, i := range tarballbasenames_s {
		parsed, err := parser.Parse(i)
		if err != nil {
			return err
		}

		tarballbasenames = append(tarballbasenames, parsed)
	}

	return self._Sort(tarballbasenames_s, tarballbasenames)
}
