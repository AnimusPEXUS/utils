package io

import "io"

type CloseDetector struct {
	CBBeforeSimple func()
	CBAfterSimple  func()
	CBBefore       func(closer io.Closer) (cancel bool, err_to_return error, force_err_to_return bool, err error)
	CBAfter        func(closer io.Closer, res error) (err_to_return error, force_err_to_return bool, err error)
	Closer         io.Closer
}

func (self *CloseDetector) Close() error {

	var err error
	var err_to_return error
	var force_err_to_return bool
	var cancel bool

	if self.CBBefore != nil {
		cancel, err_to_return, force_err_to_return, err = self.CBBefore(self.Closer)
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

	err = self.Closer.Close()

	if self.CBAfter != nil {
		err_to_return, force_err_to_return, err = self.CBAfter(self.Closer, err)
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
