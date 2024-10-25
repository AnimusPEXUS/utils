package debugging

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

type DebugToolsOptions struct {
	LogsDir       string
	ComponentName string
}

type DebugTools struct {
	options *DebugToolsOptions
}

func NewDebugTools(options *DebugToolsOptions) (*DebugTools, error) {
	self := &DebugTools{
		options: options,
	}
	return self, nil
}

func (self *DebugTools) newFile(subject string, ext string) (*os.File, error) {
	var fn string
	fn = time.Now().Format(time.RFC3339Nano)
	if self.options.ComponentName != "" {
		fn += "." + self.options.ComponentName
	}
	fn += "." + subject + "." + ext
	return os.Create(filepath.Join(self.options.LogsDir, fn))
}

func (self *DebugTools) LogRemove(filename string) error {
	return os.Remove(filepath.Join(self.options.LogsDir, filepath.Base(filename)))
}

func (self *DebugTools) LogList() ([]string, error) {
	file_list, err := ioutil.ReadDir(self.options.LogsDir)
	if err != nil {
		return nil, err
	}

	ff := make([]string, 0)
	for _, i := range file_list {
		if !i.IsDir() {
			ff = append(ff, i.Name())
		}
	}

	return ff, nil
}

func (self *DebugTools) LogSize(filename string) (int64, error) {
	s, err := os.Stat(filepath.Join(self.options.LogsDir, filepath.Base(filename)))
	if err != nil {
		return 0, err
	}
	return s.Size(), nil
}

func (self *DebugTools) LogSlice(filename string, index, size int64) ([]byte, error) {
	f, err := os.Open(filepath.Join(self.options.LogsDir, filepath.Base(filename)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, size)

	_, err = f.ReadAt(buf, index)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

func (self *DebugTools) HeapDump() (string, error) {
	f, err := self.newFile("heapdump", "bin")
	if err != nil {
		return "", err
	}
	defer f.Close()

	debug.WriteHeapDump(f.Fd())

	return f.Name(), nil
}

func (self *DebugTools) StartCPUProfile() (string, error) {
	f, err := self.newFile("profile-cpu", "bin")
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = pprof.StartCPUProfile(f)
	if err != nil {
		return f.Name(), err
	}

	return f.Name(), nil
}

func (self *DebugTools) StopCPUProfile() {
	pprof.StopCPUProfile()
	return
}

func (self *DebugTools) WriteHeapProfile() (string, error) {
	f, err := self.newFile("profile-heap", "bin")
	if err != nil {
		return "", err
	}
	defer f.Close()
	runtime.GC()
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		return f.Name(), err
	}

	return f.Name(), nil
}
