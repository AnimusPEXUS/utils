package version

import (
	"errors"
	"path"
	"sort"
	"strconv"

	"github.com/AnimusPEXUS/utils/directory"
	"github.com/AnimusPEXUS/utils/filepath"
	"github.com/AnimusPEXUS/utils/set"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

type VersionTree struct {
	d *directory.File

	tarball_name string
	//tarball_name_parser     string
	tarball_name_parser_obj types.TarballNameParserI
}

func NewVersionTree(
	tarball_name string,
	tarball_name_parser string,
) (*VersionTree, error) {
	ret := new(VersionTree)

	ret.tarball_name = tarball_name
	//ret.tarball_name_parser = tarball_name_parser

	if t, err := tarballnameparsers.Get(tarball_name_parser); err != nil {
		return nil, err
	} else {
		ret.tarball_name_parser_obj = t
	}

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

	res, err := self.tarball_name_parser_obj.ParseName(basename)
	if err != nil {
		return err
	}

	if res.Name != self.tarball_name {
		return errors.New("tarball name dismatch")
	}

	version_list := make([]uint, len(res.Version))
	copy(version_list, res.Version)

	// version_list := res.VersionStrings()

	path_part := make([]uint, 0)
	for _, i := range version_list[:len(version_list)-1] {
		path_part = append(path_part, i)
	}

	file_part := version_list[len(version_list)-1]

	directory := self.d

	if len(path_part) != 0 {
		bases := make([]string, 0)
		for _, ii := range path_part {
			ii_str := strconv.Itoa(int(ii))
			ii_f := directory.Get(ii_str)
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
					directory.MkDir(ii_str, nil)
				}
				directory = directory.Get(ii_str)
			} else {
				directory = directory.MkDir(ii_str, nil)
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
	as_name uint,
) {

	as_name_s := strconv.Itoa(int(as_name))
	if !directory.Have(as_name_s) {
		directory.MkFile(as_name_s, value)
	} else {
		d_as_name := directory.Get(as_name_s)

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
			directory = directory.Get(as_name_s)
			self.RestoreBase(directory, value, 0)
		}
	}
}

func (self *VersionTree) TruncateByVersionDepth(
	dir *directory.File,
	depth int,
) {

	if dir == nil {
		dir = self.d
	}

	lst := dir.ListDirNoSep()
	//sort.Sort(directory.FileSlice(lst))

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
		dg := dir.Get(strconv.Itoa(i))
		if dg.IsDir() {
			self.TruncateByVersionDepth(dg, depth)
		}
	}
}

func (self *VersionTree) Basenames(
	extensions_preferred_order []string,
) []string {

	bases := make([]string, 0)

	self.d.Walk(
		func(path []*directory.File, dirs, files []*directory.File) error {
			for _, i := range files {
				val := path[len(path)-1].Get(i.Name()).GetValue()

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

	return bases

}
