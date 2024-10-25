package io

import "os"

type FileCloseDetector struct {
	os.File
	CBBeforeSimple func()
	CBAfterSimple  func()
	CBBefore       func(self *FileCloseDetector) (cancel bool, err_to_return error, force_err_to_return bool, err error)
	CBAfter        func(self *FileCloseDetector, res error) (err_to_return error, force_err_to_return bool, err error)
}

func (self *FileCloseDetector) Close() error {

	var err error
	var err_to_return error
	var force_err_to_return bool
	var cancel bool

	if self.CBBefore != nil {
		cancel, err_to_return, force_err_to_return, err = self.CBBefore(self)
		if err != nil {
			panic(err)
		}

		err = err_to_return
	}

	if cancel {
		return err_to_return
	}

	if self.CBBeforeSimple != nil {
		self.CBBeforeSimple()
	}

	err = self.File.Close()

	if self.CBAfter != nil {
		err_to_return, force_err_to_return, err = self.CBAfter(self, err)
		if err != nil {
			panic(err)
		}

		if force_err_to_return {
			err = err_to_return
		}
	}

	if self.CBAfterSimple != nil {
		self.CBAfterSimple()
	}

	return err
}
