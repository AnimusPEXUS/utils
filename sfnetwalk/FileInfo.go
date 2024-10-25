package sfnetwalk

import (
	"os"
	"sort"
	"time"
)

var _ os.FileInfo = &FileInfo{}
var _ sort.Interface = OsFileInfoSort{}

type OsFileInfoSort []os.FileInfo

func (self OsFileInfoSort) Len() int           { return len(self) }
func (self OsFileInfoSort) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self OsFileInfoSort) Less(i, j int) bool { return self[i].Name() < self[j].Name() }

type FileInfo struct {
	name  string
	isdir bool
}

func (self *FileInfo) Name() string {
	return self.name
}

func (self *FileInfo) Size() int64 {
	return 0
}

func (self *FileInfo) Mode() os.FileMode {
	ret := os.FileMode(0)
	if self.isdir {
		ret |= os.ModeDir
	} else {
		ret &^= os.ModeDir
	}
	return ret
}

func (self *FileInfo) ModTime() time.Time {
	return *new(time.Time)
}

func (self *FileInfo) IsDir() bool {
	return self.isdir
}

func (self *FileInfo) Sys() interface{} {
	return nil
}

type FileInfoForMarshal struct {
	Name  string
	IsDir bool
}

func NewFileInfoForMarshal(init os.FileInfo) *FileInfoForMarshal {
	return &FileInfoForMarshal{
		Name:  init.Name(),
		IsDir: init.IsDir(),
	}
}

func (self *FileInfoForMarshal) GetFileInfo() *FileInfo {
	return &FileInfo{
		name:  self.Name,
		isdir: self.IsDir,
	}
}
