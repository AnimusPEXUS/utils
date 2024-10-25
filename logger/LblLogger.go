package logger

import (
	"bytes"
	"sync"
)

var SEP = []byte("\n")
var SEP_LEN = len(SEP)

type LineByLineWriter struct {
	logger *Logger
	err    bool
	buff   []byte
	m      *sync.Mutex
}

func NewLineByLineWriter(l *Logger, err_logger bool) *LineByLineWriter {
	ret := new(LineByLineWriter)
	ret.logger = l
	ret.err = err_logger
	ret.buff = make([]byte, 0)
	ret.m = new(sync.Mutex)
	return ret
}

func (self *LineByLineWriter) Write(p []byte) (n int, err error) {
	self.m.Lock()
	defer self.m.Unlock()
	self.buff = append(self.buff, p...)

	for {
		i := bytes.Index(self.buff, SEP)

		if i == -1 {
			break
		}

		st := string(self.buff[:i])
		self.buff = append(self.buff[:0], self.buff[i+SEP_LEN:]...)
		if !self.err {
			self.logger.Info(st)
		} else {
			self.logger.Error(st)
		}
	}

	n = len(p)
	err = nil
	return
}
