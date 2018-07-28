package pkgconfig

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/AnimusPEXUS/utils/environ"
	"github.com/AnimusPEXUS/utils/filetools"
)

type PkgConfig struct {
	executable string
}

func NewPkgConfig(searchpaths []string) (*PkgConfig, error) {

	if len(searchpaths) == 0 {
		e := environ.NewFromStrings(os.Environ())
		p := e.Get("PATH", "")

		if p == "" {
			return nil, errors.New("error determining path to search for pkg-config")
		}

		searchpaths = strings.Split(p, ":")
	}

	e, err := filetools.Which("pkg-config", searchpaths)
	if err != nil {
		return nil, err
	}

	return &PkgConfig{executable: e}, nil
}

func (self *PkgConfig) Command(args ...string) *exec.Cmd {
	return exec.Command(self.executable, args...)
}

func (self *PkgConfig) getArgsWithPrefixes(pkg_config_args []string, prefix string) ([]string, error) {
	c := self.Command(pkg_config_args...)

	var items []string

	{
		data, err := c.Output()
		if err != nil {
			return nil, err
		}

		items = strings.Split(string(data), " ")
	}

	for i := len(items) - 1; i != -1; i -= 1 {
		if !strings.HasPrefix(items[i], prefix) {
			items = append(items[:i], items[i+1:]...)
		}
	}

	prefix_len := len(prefix)

	for i := 0; i != len(items); i++ {
		items[i] = strings.Trim(items[i][prefix_len:], "\n\r\x00")
	}

	return items, nil
}

func (self *PkgConfig) GetIncludeDirs(name string) ([]string, error) {
	return self.getArgsWithPrefixes([]string{"--cflags", name}, "-I")
}

func (self *PkgConfig) GetLibNames(name string) ([]string, error) {
	return self.getArgsWithPrefixes([]string{"--libs", name}, "-l")
}
