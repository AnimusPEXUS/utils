package cache01

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/go-ini/ini"
)

const FILENAME = "config.ini"
const SECTION_NAME = "cache-cfg"
const KEY_NAME = "timeout"

type CacheDir struct {
	directory string

	settings *Settings
}

func NewCacheDir(directory string, settings *Settings) (*CacheDir, error) {
	s, err := os.Stat(directory)
	if err != nil {
		err = os.MkdirAll(directory, 0700)
		if err != nil {
			return nil, err
		}
	}

	s, err = os.Stat(directory)
	if err != nil {
		return nil, err
	}

	if !s.IsDir() {
		return nil, errors.New("must be directory")
	}
	if settings == nil {
		settings = MakeDefaultSettings()
	}
	ret := new(CacheDir)
	ret.directory = directory
	ret.settings = settings
	return ret, nil
}

func (self *CacheDir) Cache(name string, cb func() ([]byte, error)) (*Cache, error) {
	h := md5.New()
	h.Write([]byte(name))
	nn := hex.EncodeToString(h.Sum([]byte{}))
	return NewCache(self, fmt.Sprintf("%s.dat", nn), cb)
}

func (self *CacheDir) GetCfgFileName() string {
	return path.Join(self.directory, FILENAME)
}

func (self *CacheDir) GetTimeout() (time.Duration, error) {
	ret := self.settings.ListDirTimeout

	cfg_file_name := self.GetCfgFileName()

	if _, err := os.Stat(cfg_file_name); err == nil {
		tb, err := ioutil.ReadFile(cfg_file_name)
		if err != nil {
			return ret, err
		}

		i, err := ini.Load(tb)
		if err != nil {
			return ret, err
		}

		k := i.Section(SECTION_NAME).Key(KEY_NAME)
		ret = k.MustDuration(self.settings.ListDirTimeout)

	}

	return ret, nil
}

func (self *CacheDir) SetTimeout(value time.Duration) error {
	cfg_file_name := self.GetCfgFileName()

	i := ini.Empty()

	if _, err := os.Stat(cfg_file_name); err == nil {
		tb, err := ioutil.ReadFile(cfg_file_name)
		if err != nil {
			return err
		}

		i.Append(tb)
	}

	k := i.Section(SECTION_NAME).Key(KEY_NAME)
	k.SetValue(strconv.Itoa(int(value)))

	return i.SaveTo(cfg_file_name)
}
