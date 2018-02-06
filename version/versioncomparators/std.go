package versioncomparators

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

func init() {
	Index["std"] = &VersionComparatorStd{}
}

type VersionComparatorStd struct {
}

func (self *VersionComparatorStd) Compare(
	tarballbasename1, tarballbasename2 *tarballname.ParsedTarballName,
) (int, error) {

	one, err := tarballbasename1.Version.ArrInt()
	if err != nil {
		return -100, err
	}

	two, err := tarballbasename2.Version.ArrInt()
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

	sort_strings := len(tarballbasenames_s) == len(tarballbasenames)

	for i := 0; i < len(tarballbasenames)-1; i++ {
		for j := i + 1; j < len(tarballbasenames); j++ {
			pi := tarballbasenames[i]
			pj := tarballbasenames[j]

			// TODO: is this check really needed and correct?
			if pi.Name != pj.Name {
				return errors.New("by version sort name dismuch")
			}

			res, err := self.Compare(pi, pj)
			if err != nil {
				return err
			}

			if res == 1 {
				tarballbasenames[i], tarballbasenames[j] =
					tarballbasenames[j], tarballbasenames[i]
				if sort_strings {
					tarballbasenames_s[i], tarballbasenames_s[j] =
						tarballbasenames_s[j], tarballbasenames_s[i]
				}
			}
		}
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
