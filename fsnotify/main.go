// +build linux android
package fsnotify

// #include <sys/inotify.h>
import "C"

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/AnimusPEXUS/utils/worker"
	"golang.org/x/sys/unix"
)

type Event struct {
	Wd     int32
	Mask   uint32
	Cookie uint32
	Len    uint32
	Name   string
}

type Watcher struct {
	flags      int
	fd         int
	fdFile     *os.File
	filename   string
	watchdescs []uint32
	worker     *worker.Worker

	on_event_or_error_sync func(*Event, error)
}

func NewWatcher(flags int, on_event_or_error_sync func(*Event, error)) *Watcher {

	self := &Watcher{
		flags:                  flags,
		on_event_or_error_sync: on_event_or_error_sync,
	}

	self.worker = worker.New(self.readLoop)

	// NOTE: looks like it's ok to leave it empty
	self.filename = "inotify"

	return self
}

func (self *Watcher) readSuccess(e *Event) {
	self.on_event_or_error_sync(e, nil)
}

func (self *Watcher) readError(err error) {
	panic(err)
	self.on_event_or_error_sync(nil, err)
}

func (self *Watcher) GetWorker() worker.WorkerI {
	return self.worker
}

func (self *Watcher) AddWatch(path string, mask uint32) (watchdesc uint32, err error) {
	watchdesc1, err := unix.InotifyAddWatch(self.fd, path, mask)
	if err != nil {
		return
	}
	watchdesc = uint32(watchdesc1)
	self.watchdescs = append(self.watchdescs, watchdesc)
	return
}

func (self *Watcher) RmWatch(watchdesc uint32) (success int, err error) {
	success, err = unix.InotifyRmWatch(self.fd, watchdesc)
	for i := len(self.watchdescs) - 1; i != -1; i += -1 {
		if self.watchdescs[i] == watchdesc {
			self.watchdescs = append(self.watchdescs[:i], self.watchdescs[i+1:]...)
		}
	}
	return
}

func (self *Watcher) readLoop(
	set_starting func(),
	set_working func(),
	set_stopping func(),
	set_stopped func(),
	is_stop_flag func() bool,
) {
	log.Println("loop start")
	// set_starting()
	defer func() {
		log.Println("loop exiting")
		set_stopping()
		self.fdFile.Close()
		set_stopped()
		log.Println("loop exited")
	}()

	fd, err := unix.InotifyInit1(self.flags)
	// fd, err := unix.InotifyInit()
	if err != nil {
		log.Println("loop read error")
		self.readError(err)
		return
	}

	self.fd = fd
	self.fdFile = os.NewFile(uintptr(fd), self.filename)
	if self.fdFile == nil {
		self.readError(errors.New("os.NewFile returned nil" + err.Error()))
		return
	}

	stop_flag := make(chan worker.EmptyStruct)

	go func() {
		defer func() {
			go func() {
				stop_flag <- worker.EmptyStruct{}
			}()
		}()
		for {
			if is_stop_flag() {
				break
			}
			time.Sleep(time.Second)
		}
	}()

	set_working()
loop:
	for {
		select {
		case <-stop_flag:
			break loop
		case e := <-self.readNextEvent():
			self.readSuccess(e)
		}
	}

	return
}
