package cmake

import (
	"bytes"
	"go/build"
	"io"
	"os/exec"
	"path"
	"strings"
)

func _ThisDir() (string, error) {

	d, err := build.Import(
		"github.com/AnimusPEXUS/utils/cmake",
		"",
		build.FindOnly,
	)
	if err != nil {
		return "", err
	}

	return d.Dir, nil
}

type CMake struct {
	full_cmd string
	this_dir string
}

func NewCMake(
	root string,
	prefix string,
	executable string,
) (*CMake, error) {

	self := new(CMake)

	if executable == "" {
		executable = "cmake"
	} else {
		executable = path.Base(executable)
	}

	self.full_cmd = executable

	if prefix != "" {
		self.full_cmd = path.Join("/", prefix, "bin", self.full_cmd)

		if root != "" {
			self.full_cmd = path.Join("/", root, self.full_cmd)
		}
	}

	if t, err := _ThisDir(); err != nil {
		return nil, err
	} else {
		self.this_dir = t
	}

	return self, nil
}

func (self *CMake) GetExecutable() string {
	return self.full_cmd
}

func (self *CMake) GetCMAKE_ROOT() (string, error) {

	params := []string{"-P", path.Join(self.this_dir, "root.cmake")}

	c := exec.Command(self.full_cmd, params...)
	c.Dir = self.this_dir

	errp, err := c.StderrPipe()
	buff := &bytes.Buffer{}

	err = c.Start()
	if err != nil {
		return "", err
	}

	ce := make(chan bool)

	go func(errp io.ReadCloser, buff *bytes.Buffer) {
		_, err := io.Copy(buff, errp)
		if err != nil {
			panic(err)
		}
		ce <- true
	}(errp, buff)

	<-ce

	out := string(buff.Bytes())
	out = strings.Trim(out, "\n\r\x00")

	return out, nil
}
