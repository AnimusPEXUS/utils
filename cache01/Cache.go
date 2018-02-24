package cache01

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Cache struct {
	dir      *CacheDir
	filename string
	cb       func() ([]byte, error)
}

func NewCache(
	dir *CacheDir,
	filename string,
	// callback to get fresh value
	cb func() ([]byte, error),
) (*Cache, error) {
	ret := new(Cache)
	ret.dir = dir
	ret.filename = filename
	ret.cb = cb

	if s, err := os.Stat(ret.dir.directory); err != nil {
		return nil, err
	} else {
		if !s.IsDir() {
			err = os.MkdirAll(ret.dir.directory, 0700)
			if err != nil {
				return nil, errors.New("not a directory")
			}
		}
	}

	return ret, nil
}

func (self *Cache) GetFileName() string {
	return path.Join(self.dir.directory, self.filename)
}

func (self *Cache) GetValue() ([]byte, error) {
	// TODO: should be added file integrety check
	filename := self.GetFileName()

	get_new_data := false

	s, err := os.Stat(filename)
	if err != nil {
		get_new_data = true
	} else {
		mt := s.ModTime().UTC()
		nt := time.Now().UTC()

		dif := nt.Sub(mt)

		if dif < 0 {
			get_new_data = true
		} else {
			to, err := self.dir.GetTimeout()
			if err != nil {
				return nil, err
			}
			if dif > to {
				get_new_data = true
			}
		}
	}

	if get_new_data {
		data, err := self.cb()
		if err != nil {
			return nil, err
		}
		err = self.SetValue(data)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

}

func (self *Cache) SetValue(value []byte) error {
	filename := self.GetFileName()
	err := ioutil.WriteFile(filename, value, 0700)
	if err != nil {
		return err
	}
	return nil
}
