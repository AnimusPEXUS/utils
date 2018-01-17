package directory

import "fmt"

// TODO: probably mutexes should be added

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

func (self *File) GetThisDirPath() []*File {

	if !self.is_directory {
		panic("can't be used for not dir")
	}

	ret := make([]*File, 0)

	p := self

	for p != nil {

		ret = append(ret, nil)
		copy(ret[1:], ret[0:])
		ret[0] = p

		p = p.Parent()

	}

	return ret
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
) *File {

	if debug {
		fmt.Println("MkDir", name, value)
	}

	if !self.is_directory {
		panic("can't be used for not dir")
	}

	self.Delete(name)

	f := NewFile(self, name, true, value)
	self.directory = append(self.directory, f)

	if debug {
		fmt.Println("   ret", f)
	}
	return f
}

func (self *File) MkFile(
	name string,
	value interface{},
) *File {

	if debug {
		fmt.Println("MkFile", name, value)
	}

	if !self.is_directory {
		panic("can't be used for not dir")
	}

	self.Delete(name)

	f := NewFile(self, name, false, value)
	self.directory = append(self.directory, f)

	if debug {
		fmt.Println("   ret", f)
	}
	return f
}

func (self *File) Delete(name string) {

	if debug {
		fmt.Println("Delete", name)
	}

	if !self.is_directory {
		panic("can't be used for not dir")
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
	return
}

func (self *File) ListDirNoSep() []*File {

	if !self.is_directory {
		panic("can't be used for not dir")
	}

	ret := make([]*File, len(self.directory))
	copy(ret, self.directory)

	return ret
}

func (self *File) ListDir() ([]*File, []*File) {

	if !self.is_directory {
		panic("can't be used for not dir")
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

	return dirs, files
}

func (self *File) Walk(
	target func(path []*File, dirs, files []*File) error,
) error {

	if !self.is_directory {
		panic("can't be used for not dir")
	}

	dirs, files := self.ListDir()

	path := self.GetThisDirPath()

	err := target(path, dirs, files)

	for _, i := range dirs {
		err = i.Walk(target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *File) Have(name string) bool {

	if !self.is_directory {
		panic("can't be used for not dir")
	}

	return self.Get(name) != nil
}

func (self *File) Get(name string) *File {

	if !self.is_directory {
		panic("can't be used for not dir")
	}

	for _, i := range self.directory {
		if i.Name() == name {
			return i
		}
	}

	return nil
}

func (self *File) GetByPath(
	pth []string,
	create_if_not_exists bool,
	last_element_is_file bool,
	value interface{},
) *File {

	var ret *File = nil

	if len(pth) != 0 {
		name := pth[0]
		if name == "" {
			ret = self.GetRoot()
		} else if t := self.Get(name); t != nil {
			ret = t
		} else {
			if create_if_not_exists {
				if len(pth) == 1 && last_element_is_file {
					ret = self.MkFile(name, value)
				} else {
					ret = self.MkDir(name, value)
				}
			} else {
				ret = nil
			}
		}
	} else {
		ret = self
	}

	if ret != nil && len(pth) != 0 {

		rem_pth := make([]string, len(pth)-1)
		copy(rem_pth, pth[1:])
		if len(rem_pth) == 0 {

		} else {
			ret = ret.GetByPath(
				rem_pth,
				create_if_not_exists,
				last_element_is_file,
				value,
			)
		}

	}

	return ret

}

func (self *File) SetValue(value interface{}) {
	self.value = value
}

func (self *File) GetValue() interface{} {
	return self.value
}
