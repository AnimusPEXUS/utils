package gosignal

import (
	"fmt"
	"sync"
)

type Signal1 struct {
	listeners     map[uint64]ListenerFunction
	emition_mutex *sync.Mutex
	counter       uint64
}

func NewSignal1() *Signal1 {

	ret := new(Signal1)
	ret.emition_mutex = new(sync.Mutex)
	ret.listeners = make(map[uint64]ListenerFunction)
	ret.counter = 0

	return ret
}

func (self *Signal1) Connect(f ListenerFunction) *Connector1 {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	if DEBUG {
		fmt.Println("signal", self, "connecting to", f)
	}

	self.listeners[self.counter] = f

	ret := new(Connector1)
	ret.id = self.counter
	ret.s = self
	if DEBUG {
		fmt.Println("   ", "created", ret, "with id", ret.id)
	}

	self.counter++

	return ret
}

func (self *Signal1) Disconnect(id uint64) {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	if DEBUG {
		fmt.Println("disconnecting object with id", id)
	}

	delete(self.listeners, id)

}

func (self *Signal1) Emit(data interface{}) {

	defer self.emition_mutex.Unlock()
	self.emition_mutex.Lock()

	for _, i := range self.listeners {
		if DEBUG {
			fmt.Println("emiting to", i)
		}

		go i(data)
	}
}

type Connector1 struct {
	s  *Signal1
	id uint64
}

func (self *Connector1) Disconnect() {
	self.s.Disconnect(self.id)
}

type Connector1Pool struct {
	lst []*Connector1
}

func (self *Connector1Pool) Disconnect() {
	for _, i := range self.lst {
		i.Disconnect()
	}
}
