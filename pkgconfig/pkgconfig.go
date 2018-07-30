package pkgconfig

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/AnimusPEXUS/utils/environ"
	"github.com/AnimusPEXUS/utils/filetools"
	"github.com/AnimusPEXUS/utils/set"
)

// TODO: I don't like this.. I already have such list somethere
var POSSIBLE_LIBDIR_NAMES = []string{"lib", "lib64"}

type PkgConfig struct {
	executable          string
	pkg_config_path_env []string
}

func NewPkgConfig(
	path_env []string,
	pkg_config_path_env []string,
) (*PkgConfig, error) {

	if len(path_env) == 0 {
		e := environ.NewFromStrings(os.Environ())
		p := e.Get("PATH", "")

		if p == "" {
			return nil, errors.New("error determining path to search for pkg-config")
		}

		path_env = strings.Split(p, ":")
	}

	if len(pkg_config_path_env) == 0 {
		r := make([]string, 0)

		prefixes := make([]string, 0)

		{
			prefixes_s := set.NewSetString()

			for _, i := range path_env {
				prefixes_s.AddStrings(path.Base(i))
			}

			prefixes = prefixes_s.ListStrings()
		}

		dirs := make([]string, 0)
		dirs = append(dirs, POSSIBLE_LIBDIR_NAMES...)
		dirs = append(dirs, "share")

		for _, i := range prefixes {
			for _, j := range dirs {
				pth := path.Join(i, j, "pkgconfig")
				if s, err := os.Stat(pth); err != nil {
					if !os.IsNotExist(err) {
						return nil, err
					} else {
						continue
					}
				} else {
					if !s.IsDir() {
						return nil, errors.New(pth + " isn't a directory")
					}
				}
				r = append(r, pth)
			}
		}

		pkg_config_path_env = r

	}

	e, err := filetools.Which("pkg-config", path_env)
	if err != nil {
		return nil, err
	}

	ret := &PkgConfig{
		executable:          e,
		pkg_config_path_env: pkg_config_path_env,
	}

	return ret, nil
}

func (self *PkgConfig) Command(args ...string) *exec.Cmd {

	env := environ.NewFromStrings(os.Environ())
	env.Set(
		"PKG_CONFIG_PATH",
		strings.Join(
			self.pkg_config_path_env,
			string(os.PathListSeparator),
		),
	)

	c := exec.Command(self.executable, args...)
	c.Env = env.Strings()
	return c
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
