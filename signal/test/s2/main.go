package main

import (
	"log"
	"runtime"
	"time"

	"github.com/AnimusPEXUS/utils/signal"
)

type A struct {
	Signal0 *signal.Signal2
}

func NewA() *A {
	self := &A{
		Signal0: signal.NewSignal2(),
	}
	go self.thread()
	return self
}

func (self *A) thread() {
	for {
		time.Sleep(time.Second)
		log.Print("emiting")
		self.Signal0.Emit(123)
	}
}

type B struct {
}

func (self *B) handler(data interface{}) {
	log.Print("handler worked")
}

func main() {

	a := NewA()
	b := &B{}

	time.Sleep(5 * time.Second)
	connector := a.Signal0.Connect(b.handler)
	connector.Use()
	time.Sleep(5 * time.Second)
	connector = nil
	runtime.GC()
	time.Sleep(5 * time.Second)

}
