package directory

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

// TODO: probably mutexes should be added

var ERR_CANT_BE_USED_ON_NON_DIRECTORY = errors.New("can't be used for not dir")

const debug = false

type FileSlice []*File

func (self FileSlice) Len() int           { return len(self) }
func (self FileSlice) Swap(i, j int)      { self[j], self[i] = self[i], self[j] }
func (self FileSlice) Less(i, j int) bool { return self[i].Name() < self[j].Name() }

type File struct {
	name   string
	parent *File

	is_directory bool

	directory []*File
	value     interface{}
}

func NewTree() *File {
	return NewFile(nil, "", true, nil)
}

func NewFile(
	parent *File,
	name string,
	is_directory bool,
	value interface{},
) *File {

	if debug {
		fmt.Println("NewFile", parent, name, is_directory, value)
	}

	ret := new(File)

	ret.is_directory = is_directory

	if is_directory {
		ret.directory = make([]*File, 0)
	}
	ret.value = value
	ret.parent = parent
	ret.name = name

	if debug {
		fmt.Println("   ret", ret)
	}

	return ret
}

// func (self *File) _GetMutex() *sync.RWMutex {
// 	return self.GetRoot().m
// }

func (self *File) Name() string {
	return self.name
}

func (self *File) IsDir() bool {
	return self.is_directory
}

func (self *File) Parent() *File {
	return self.parent
}

func (self *File) PathString() (string, error) {
	ret_lst := make([]string, 0)
	p, err := self.Path()
	if err != nil {
		return "", err
	}
	for _, i := range p {
		ret_lst = append(ret_lst, i.Name())
	}
	return strings.Join(ret_lst, "/"), nil
}

func (self *File) Path() ([]*File, error) {
	ret := make([]*File, 0)

	p := self

	for {
		ret = append([]*File{p}, ret...)

		p = p.parent

		if p == nil {
			break
		}
	}

	return ret, nil
}

func (self *File) GetRoot() *File {

	ret := self

	for {
		p := ret.Parent()
		if p == nil {
			break
		}
		ret = p
	}

	return ret
}

func (self *File) MkDir(
	name string,
	value interface{},
) (*File, error) {

	if debug {
		fmt.Println("MkDir", name, value)
	}

	if !self.is_directory {
		return nil, ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	self.Delete(name)

	f := NewFile(self, name, true, value)
	self.directory = append(self.directory, f)

	if debug {
		fmt.Println("   ret", f)
	}
	return f, nil
}

func (self *File) MkFile(
	name string,
	value interface{},
) (*File, error) {

	if debug {
		fmt.Println("MkFile", name, value)
	}

	if !self.is_directory {
		return nil, ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	self.Delete(name)

	f := NewFile(self, name, false, value)
	self.directory = append(self.directory, f)

	if debug {
		fmt.Println("   ret", f)
	}
	return f, nil
}

func (self *File) Delete(name string) error {

	if debug {
		fmt.Println("Delete", name)
	}

	if !self.is_directory {
		return ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	index := -1
	for k, v := range self.directory {
		if v.Name() == name {
			index = k
			break
		}
	}
	if index != -1 {
		self.directory = append(self.directory[:index], self.directory[index+1:]...)
	}
	return nil
}

func (self *File) ListDirNoSep() ([]*File, error) {

	if !self.is_directory {
		return []*File{}, ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	ret := make([]*File, len(self.directory))
	copy(ret, self.directory)

	return ret, nil
}

func (self *File) ListDir() ([]*File, []*File, error) {

	if !self.is_directory {
		return []*File{}, []*File{}, ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	dirs := make([]*File, 0)
	files := make([]*File, 0)

	for _, i := range self.directory {
		if i.IsDir() {
			dirs = append(dirs, i)
		} else {
			files = append(files, i)
		}
	}

	return dirs, files, nil
}

func (self *File) Walk(
	target func(path, dirs, files []*File) error,
) error {

	if !self.is_directory {
		return ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	dirs, files, err := self.ListDir()
	if err != nil {
		return err
	}

	path, err := self.Path()
	if err != nil {
		return err
	}

	err = target(path, dirs, files)

	for _, i := range dirs {
		err = i.Walk(target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *File) Have(name string) (bool, error) {

	// already in .Get() method
	// if !self.is_directory {
	// 	return false, ERR_CANT_BE_USED_ON_NON_DIRECTORY
	// }

	get_res, err := self.Get(name, false)
	if err != nil {
		return false, err
	}

	return get_res != nil, nil
}

func (self *File) Get(name string, not_found_is_error bool) (*File, error) {

	if !self.is_directory {
		return nil, ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	for _, i := range self.directory {
		if i.Name() == name {
			return i, nil
		}
	}

	if not_found_is_error {
		return nil, errors.New("not found")
	}

	return nil, nil
}

// func (self *File) GetByPath(
// 	pth []string,
// 	create_if_not_exists bool,
// 	last_element_is_file bool,
// 	value interface{},
// ) *File {
//
// 	var ret *File = nil
//
// 	if len(pth) != 0 {
// 		name := pth[0]
// 		if name == "" {
// 			ret = self.GetRoot()
// 		} else if t := self.Get(name); t != nil {
// 			ret = t
// 		} else {
// 			if create_if_not_exists {
// 				if len(pth) == 1 && last_element_is_file {
// 					ret = self.MkFile(name, value)
// 				} else {
// 					ret = self.MkDir(name, value)
// 				}
// 			} else {
// 				ret = nil
// 			}
// 		}
// 	} else {
// 		ret = self
// 	}
//
// 	if ret != nil && len(pth) != 0 {
//
// 		rem_pth := make([]string, len(pth)-1)
// 		copy(rem_pth, pth[1:])
// 		if len(rem_pth) == 0 {
//
// 		} else {
// 			ret = ret.GetByPath(
// 				rem_pth,
// 				create_if_not_exists,
// 				last_element_is_file,
// 				value,
// 			)
// 		}
//
// 	}
//
// 	return ret
// }

func (self *File) GetByPath(
	pth []string,
	create_if_not_exists bool,
	last_element_is_file bool,
	value interface{},
) (*File, error) {

	if self.parent != nil {
		return nil, errors.New("this method works only with root File")
	}

	if !self.is_directory {
		return nil, ERR_CANT_BE_USED_ON_NON_DIRECTORY
	}

	if len(pth) == 0 {
		return self, nil
	}

	var ret *File

	dirs := pth[0 : len(pth)-1]
	if !last_element_is_file {
		dirs = append(dirs, pth[len(pth)-1])
	}

	var d_track *File = self

	for _, i := range dirs {
		// var get_res *File
		get_res, err := d_track.Get(i, false)
		if err != nil {
			return nil, err
		}
		if get_res == nil {
			if create_if_not_exists {
				get_res, err = d_track.MkDir(i, nil)
				if err != nil {
					return nil, err
				}
				d_track = get_res
			} else {
				return nil, errors.New("some path element is not exists")
			}
		} else {
			if !get_res.IsDir() {
				return nil, errors.New(
					"path to last element exists, but one of elements isn't dir",
				)
			}
			d_track = get_res
		}
	}

	if last_element_is_file {
		n := pth[len(pth)-1]
		get_res, err := d_track.Get(n, false)
		if err != nil {
			return nil, err
		}
		if get_res == nil {
			if create_if_not_exists {
				get_res, err = d_track.MkFile(n, value)
				if err != nil {
					return nil, err
				}
				ret = get_res
			}
		} else {
			ret = get_res
		}
	} else {
		ret = d_track
	}

	return ret, nil
}

func (self *File) SetValue(value interface{}) {
	self.value = value
}

func (self *File) GetValue() interface{} {
	return self.value
}

func (self *File) FindFile(pattern string) ([]string, error) {
	// if self.parent == nil {
	// 	return []string{}, errors.New("this method allowed only on root")
	// }
	ret := make([]string, 0)
	err := self.Walk(
		func(pth, dirs, files []*File) error {
			for _, i := range files {
				path_str, err := pth[len(pth)-1].PathString()
				if err != nil {
					return err
				}
				i_n := i.Name()
				m, err := path.Match(pattern, i_n)
				if err != nil {
					return err
				}
				if m {
					path_str = path_str + "/" + i_n
					ret = append(ret, path_str)
				}
			}
			return nil
		},
	)
	if err != nil {
		return []string{}, err
	}
	return ret, nil
}

func (self *File) TreeString() (string, error) {
	ret := ""
	err := self.Walk(
		func(path, dirs, files []*File) error {
			p, err := path[len(path)-1].PathString()
			if err != nil {
				return err
			}
			ret += fmt.Sprintln("path", p)
			for ii, i := range dirs {
				ret += fmt.Sprintln("    <dir>", ii, i.Name())
			}
			for ii, i := range files {
				ret += fmt.Sprintln("    ", ii, i.Name())
			}
			return nil
		},
	)
	if err != nil {
		return "", err
	}
	return ret, nil
}
