package tarballversion

import (
	"errors"
	"path"
	"regexp"
	"sort"
	"strconv"

	"github.com/AnimusPEXUS/utils/directory"
	"github.com/AnimusPEXUS/utils/filepath"
	"github.com/AnimusPEXUS/utils/set"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
	types2 "github.com/AnimusPEXUS/utils/tarballversion/versioncomparators/types"
)

type VersionTree struct {
	d *directory.File

	tarball_name           string
	tarball_name_is_regexp bool
	tarball_name_parser    types.TarballNameParserI
	comparator             types2.VersionComparatorI
}

func NewVersionTree(
	tarball_name string,
	tarball_name_is_regexp bool,
	tarball_name_parser types.TarballNameParserI,
	comparator types2.VersionComparatorI,
) (*VersionTree, error) {
	ret := new(VersionTree)

	ret.tarball_name = tarball_name
	ret.tarball_name_is_regexp = tarball_name_is_regexp
	ret.tarball_name_parser = tarball_name_parser
	ret.comparator = comparator

	// if t, err := tarballnameparsers.Get(tarball_name_parser); err != nil {
	// 	return nil, err
	// } else {
	// 	ret.tarball_name_parser_obj = tarball_name_parser
	// }

	ret.d = directory.NewFile(
		nil,
		"",
		true,
		nil,
	)

	return ret, nil
}

func (self *VersionTree) Add(basename string) error {

	basename = path.Base(basename)

	res, err := self.tarball_name_parser.Parse(basename)
	if err != nil {
		return err
	}

	var match bool

	if self.tarball_name_is_regexp {

		if m, err := regexp.MatchString(self.tarball_name, res.Name); err != nil {
			return err
		} else {
			match = m
		}
	} else {
		match = self.tarball_name == res.Name
	}

	if !match {
		return errors.New("tarball name dismatch")
	}

	version_list, err := self.comparator.RenderNumericalVersion(res)
	if err != nil {
		return err
	}

	// version_list, err := res.Version.ArrInt()
	// if err != nil {
	// 	return err
	// }

	path_part := make([]int, 0)
	for _, i := range version_list[:len(version_list)-1] {
		path_part = append(path_part, i)
	}

	file_part := version_list[len(version_list)-1]

	directory := self.d

	if len(path_part) != 0 {
		bases := make([]string, 0)
		for _, ii := range path_part {
			ii_str := strconv.Itoa(int(ii))
			ii_f, err := directory.Get(ii_str, false)
			if err != nil {
				return err
			}
			if ii_f != nil {
				if !ii_f.IsDir() {
					ii_f_v := ii_f.GetValue()
					switch ii_f_v.(type) {
					case []string:
						bases = append(bases, ii_f.GetValue().([]string)...)
					default:
						panic("programming error")
					}

					//directory.Delete(ii_str)
					_, err = directory.MkDir(ii_str, nil)
					if err != nil {
						return err
					}
				}
				directory, err = directory.Get(ii_str, false)
				if err != nil {
					return err
				}
			} else {
				directory, err = directory.MkDir(ii_str, nil)
				if err != nil {
					return err
				}
			}
		}
		for _, i := range bases {
			self.RestoreBase(directory, []string{i}, 0)
		}
	}

	self.RestoreBase(directory, []string{basename}, file_part)

	return nil

}

func (self *VersionTree) RestoreBase(
	directory *directory.File,
	value []string,
	as_name int,
) error {

	as_name_s := strconv.Itoa(int(as_name))
	h, err := directory.Have(as_name_s)
	if err != nil {
		return err
	}
	if !h {
		_, err = directory.MkFile(as_name_s, value)
		if err != nil {
			return err
		}
	} else {
		d_as_name, err := directory.Get(as_name_s, true)
		if err != nil {
			return err
		}

		if !d_as_name.IsDir() {
			fv := d_as_name.GetValue()

			switch fv.(type) {
			case string:
				d_as_name.SetValue([]string{fv.(string)})
				fv = d_as_name.GetValue()
			}

			fv = append(fv.([]string), value...)

			s := set.NewSetString()

			for _, i := range fv.([]string) {
				s.Add(i)
			}

			d_as_name.SetValue(s.ListStrings())

		} else {
			directory, err = directory.Get(as_name_s, true)
			if err != nil {
				return err
			}
			self.RestoreBase(directory, value, 0)
		}
	}
	return nil
}

func (self *VersionTree) TruncateByVersionDepth(
	dir *directory.File,
	depth int,
) error {

	if depth < 0 {
		return nil
	}

	if dir == nil {
		dir = self.d
	}

	lst, err := dir.ListDirNoSep()
	if err != nil {
		return err
	}
	// sort.Sort(directory.FileSlice(lst))

	inames := make([]int, 0)

	for _, i := range lst {
		i_n_t, err := strconv.Atoi(i.Name())
		if err != nil {
			panic("programming error: version file \\names\\ must be decimals")
		}
		inames = append(inames, i_n_t)
	}

	sort.Ints(inames)

	inames_to_delete := make([]int, 0)
	inames_to_work_with := make([]int, 0)

	count := 0

	for i := len(inames) - 1; i != -1; i-- {
		if count < depth {
			inames_to_work_with = append(inames_to_work_with, inames[i])
		} else {
			inames_to_delete = append(inames_to_delete, inames[i])
		}
		count++
	}

	for _, i := range inames_to_delete {
		dir.Delete(strconv.Itoa(i))
	}

	for _, i := range inames_to_work_with {
		dg, err := dir.Get(strconv.Itoa(i), true)
		if err != nil {
			return err
		}
		if dg.IsDir() {
			self.TruncateByVersionDepth(dg, depth)
		}
	}

	return nil
}

func (self *VersionTree) Basenames(
	extensions_preferred_order []string,
) ([]string, error) {

	bases := make([]string, 0)

	err := self.d.Walk(
		func(path, dirs, files []*directory.File) error {
			val_t := self.d.GetRoot()

			if len(path) != 0 {
				val_t = path[len(path)-1]
			}

			for _, i := range files {
				val_t2, err := val_t.Get(i.Name(), true)
				if err != nil {
					return err
				}
				val := val_t2.GetValue()

				res := filepath.SelectByPreferredExtension(
					val.([]string),
					extensions_preferred_order,
				)

				if res != "" {
					bases = append(bases, res)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return bases, nil

}

func (self *VersionTree) TreeString() (string, error) {
	return self.d.TreeString()
}
