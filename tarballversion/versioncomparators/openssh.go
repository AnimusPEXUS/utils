package versioncomparators

import (
	"errors"
	"strconv"

	"github.com/AnimusPEXUS/utils/sort"
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
	"github.com/AnimusPEXUS/utils/versionorstatus"
)

func init() {
	Index["openssh"] = &VersionComparatorOpenSSH{}
}

type VersionComparatorOpenSSH struct {
}

func (self *VersionComparatorOpenSSH) RenderNumericalVersion(
	tarballbasename *tarballname.ParsedTarballName,
) (
	[]int, error,
) {

	ret, err := tarballbasename.Version.ArrInt()
	if err != nil {
		return nil, err
	}

	len_arr := len(tarballbasename.Status.Arr)

	if !(len_arr == 0 || len_arr == 2) {
		return nil, errors.New("invalid number of elements in status")
	}

	if len_arr == 0 {
		return ret, nil
	}

	p_num := tarballbasename.Status.Arr[1]
	p_num_i, err := strconv.Atoi(p_num)
	if err != nil {
		return nil, err
	}
	ret = append(ret, p_num_i)

	return ret, nil
}

func (self *VersionComparatorOpenSSH) Compare(
	tarballbasename1, tarballbasename2 *tarballname.ParsedTarballName,
) (int, error) {
	return Index["std"].Compare(tarballbasename1, tarballbasename2)
}

func (self *VersionComparatorOpenSSH) _Sort(
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

func (self *VersionComparatorOpenSSH) Sort(
	tarballbasenames []*tarballname.ParsedTarballName,
) error {
	return self._Sort([]string{}, tarballbasenames)
}

func (self *VersionComparatorOpenSSH) SortStrings(
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
