package io

import "io"

// this is experimental and not well tested yet

// var _ io.WriterAt = (interface{}(CloseDetector{})).(io.WriterAt)

type CloseDetector struct {
	Object         interface{}
	CBBeforeSimple func()
	CBAfterSimple  func()
	CBBefore       func(self *CloseDetector) (cancel bool, err_to_return error, force_err_to_return bool, err error)
	CBAfter        func(self *CloseDetector, res error) (err_to_return error, force_err_to_return bool, err error)
}

func (self *CloseDetector) Close() error {

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

	{
		closer, ok := self.Object.(io.Closer)
		if ok {
			err = closer.Close()
		}
	}

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
