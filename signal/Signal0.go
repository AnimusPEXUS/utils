package signal

import (
	"log"
	"runtime"
	"sync"
	"unsafe"
)

/*
	tryings to use runtime module to track for connected method unavailability
	doesn't works too well. use Signal1
*/
type Signal0 struct {
	listeners     map[uint64]uintptr
	emition_mutex *sync.Mutex
	counter       uint64
}

func NewSignal0() *Signal0 {

	ret := new(Signal0)
	ret.emition_mutex = new(sync.Mutex)
	ret.listeners = make(map[uint64]uintptr)
	ret.counter = 0

	return ret
}

func (self *Signal0) Connect(f ListenerFunction) *Connector0 {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	if DEBUG {
		log.Println("signal", self, "connecting to", f)
	}

	self.listeners[self.counter] = uintptr(unsafe.Pointer(&f))

	ret := new(Connector0)
	ret.id = self.counter
	ret.s = self
	if DEBUG {
		log.Println("   ", "created", ret, "with id", ret.id)
	}

	runtime.SetFinalizer(ret, self.connector_finalizer)

	self.counter++

	return ret
}

func (self *Signal0) Disconnect(id uint64) {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	if DEBUG {
		log.Println("disconnecting object with id", id)
	}

	delete(self.listeners, id)

}

func (self *Signal0) connector_finalizer(obj *Connector0) {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	if DEBUG {
		log.Println("finalizing", obj, "with id", obj.id)
	}

	delete(self.listeners, obj.id)
}

func (self *Signal0) Emit(data interface{}) {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	for _, i := range self.listeners {
		if DEBUG {
			log.Println("emiting to", i)
		}

		p := (unsafe.Pointer)(i)
		f := *(*ListenerFunction)(p)

		go f(data)
	}
}

type Connector0 struct {
	s  *Signal0
	id uint64
}

func (self *Connector0) Disconnect() {
	self.s.Disconnect(self.id)
}

type Connector0Pool struct {
	lst []*Connector1
}

func (self *Connector0Pool) Disconnect() {
	for _, i := range self.lst {
		i.Disconnect()
	}
}
