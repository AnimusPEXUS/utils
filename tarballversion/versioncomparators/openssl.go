package versioncomparators

import (
	"errors"
	"regexp"
	"strings"

	"github.com/AnimusPEXUS/utils/sort"
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
	"github.com/AnimusPEXUS/utils/versionorstatus"
)

func init() {
	Index["openssl"] = &VersionComparatorOpenSSL{}
}

type VersionComparatorOpenSSL struct {
}

func (self *VersionComparatorOpenSSL) RenderNumericalVersion(
	tarballbasename *tarballname.ParsedTarballName,
) (
	[]int, error,
) {

	len_arr := len(tarballbasename.Status.Arr)

	version_to_add, err := tarballbasename.Version.ArrInt()
	if err != nil {
		return nil, err
	}

	if len_arr > 1 {
		return nil, errors.New("unacceptable status to parse as openssl release")
	}

	if len_arr == 1 {

		if ok, err := regexp.MatchString(`^[a-z]{1,3}$`, tarballbasename.Status.Arr[0]); err != nil {
			return nil, err
		} else {
			if !ok {
				return nil, errors.New("unacceptable status to parse as openssl release")
			}
		}

		letter_versions_int := make([]int, 0)

		{
			stat_arr0 := tarballbasename.Status.Arr[0]
			splitted_stat := strings.Split(stat_arr0, "")

			for _, i := range splitted_stat {
				ii := int(byte(i[0]) - 96)
				letter_versions_int = append(letter_versions_int, ii)
			}
		}

		version_to_add = append(version_to_add, letter_versions_int...)
	}

	return version_to_add, nil
}

func (self *VersionComparatorOpenSSL) Compare(
	tarballbasename1, tarballbasename2 *tarballname.ParsedTarballName,
) (int, error) {
	return Index["std"].Compare(tarballbasename1, tarballbasename2)
}

func (self *VersionComparatorOpenSSL) _Sort(
	tarballbasenames_s []string,
	tarballbasenames []*tarballname.ParsedTarballName,
) error {

	basenames_versions := make([][]int, 0)

	for _, i := range tarballbasenames {
		version_to_add, err := self.RenderNumericalVersion(i)
		if err != nil {
			return err
		}
		basenames_versions = append(basenames_versions, version_to_add)
	}

	what_to_sort := []interface{}{
		tarballbasenames,
		basenames_versions,
	}

	if len(tarballbasenames_s) == len(tarballbasenames) {
		what_to_sort = append(what_to_sort, tarballbasenames_s)
	}

	err := sort.Sort(
		what_to_sort,
		1,
		func(
			i interface{},
			j interface{},
		) (int, error) {
			pi := &tarballname.ParsedTarballName{
				Name:    "aaa",
				Version: versionorstatus.NewParsedVersionFromArrInt(i.([]int)),
			}

			pj := &tarballname.ParsedTarballName{
				Name:    "aaa",
				Version: versionorstatus.NewParsedVersionFromArrInt(j.([]int)),
			}

			// TODO: is this check really needed and correct?
			if pi.Name != pj.Name {
				return -100, errors.New("by version sort name dismuch")
			}

			res, err := self.Compare(pi, pj)
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

func (self *VersionComparatorOpenSSL) Sort(
	tarballbasenames []*tarballname.ParsedTarballName,
) error {
	return self._Sort([]string{}, tarballbasenames)
}

func (self *VersionComparatorOpenSSL) SortStrings(
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
