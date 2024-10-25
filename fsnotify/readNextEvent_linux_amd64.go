package fsnotify

import (
	"encoding/binary"
	"io"

	"golang.org/x/sys/unix"
)

func (self *Watcher) readNextEvent() <-chan *Event {
	c := make(chan *Event)
	go func() {
		x := &Event{}
		defer func() { c <- x }()

		fds := &unix.FdSet{
			Bits: [16]int64{int64(self.fdFile.Fd())},
		}

		unix.Select(2, fds, nil, nil, nil)
		// b4 := make([]byte, 4)

		// rc, err := unix.Read(self.fd, b4)
		// rc, err := io.ReadFull(self.fdFile, b4)

		// log.Println("rc, err ", rc, err)

		err := binary.Read(self.fdFile, binary.LittleEndian, &x.Wd)
		if err != nil {
			self.readError(err)
			return
		}

		err = binary.Read(self.fdFile, binary.LittleEndian, &x.Mask)
		if err != nil {
			self.readError(err)
			return
		}

		err = binary.Read(self.fdFile, binary.LittleEndian, &x.Cookie)
		if err != nil {
			self.readError(err)
			return
		}

		err = binary.Read(self.fdFile, binary.LittleEndian, &x.Len)
		if err != nil {
			self.readError(err)
			return
		}

		name := make([]byte, x.Len)

		_, err = io.ReadFull(self.fdFile, name)
		if err != nil {
			self.readError(err)
			return
		}

		x.Name = string(name)

	}()
	return c
}
