package signal

import (
	"log"
	"runtime"
	"sync"
	"unsafe"
)

type Signal2 struct {
	listeners     []uintptr
	emition_mutex *sync.Mutex
}

func NewSignal2() *Signal2 {

	ret := new(Signal2)
	ret.emition_mutex = new(sync.Mutex)
	ret.listeners = make([]uintptr, 0)

	return ret
}

func (self *Signal2) Connect(f ListenerFunction) *Connector2 {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	if DEBUG {
		log.Println("signal", self, "connecting to", f)
	}

	ptr := uintptr(unsafe.Pointer(&f))

	for _, i := range self.listeners {
		if i == ptr {
			return nil
		}
	}

	self.listeners = append(self.listeners, ptr)

	ret := new(Connector2)
	ret.id = ptr
	ret.s = self
	ret.f = f
	if DEBUG {
		log.Println("   ", "created", ret, "with id", ret.id)
	}

	runtime.SetFinalizer(ret, self.connector_finalizer)

	return ret
}

func (self *Signal2) disconnect(id uintptr) {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	if DEBUG {
		log.Println("disconnecting object with id", id)
	}

	for i := len(self.listeners) - 1; i != -1; i += -1 {
		v := self.listeners[i]
		if v == id {
			self.listeners = append(self.listeners[:i], self.listeners[i+1:]...)
		}
	}

}

func (self *Signal2) connector_finalizer(obj *Connector2) {

	if DEBUG {
		log.Println("finalizing", obj, "with id", obj.id)
	}

	self.disconnect(obj.id)
}

func (self *Signal2) Emit(data interface{}) {

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

type Connector2 struct {
	s  *Signal2
	id uintptr
	f  ListenerFunction
}

func (self *Connector2) Disconnect() {
	self.s.disconnect(self.id)
}

func (self *Connector2) Use() {}

type Connector2Pool struct {
	lst []*Connector1
}

func (self *Connector2Pool) Disconnect() {
	for _, i := range self.lst {
		i.Disconnect()
	}
}
