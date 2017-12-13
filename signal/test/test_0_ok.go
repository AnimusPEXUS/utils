package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/AnimusPEXUS/gosignal"
)

type A struct {
	c *gosignal.ConnectorUnsafe
}

func (self *A) Emited(data interface{}) {
	fmt.Println(self.Emited, "got signal", data)
}

func main() {
	s := gosignal.NewSignalUnsafe()

	{
		b := make([]*A, 0)

		for i := 0; i != 10; i++ {
			o := &A{}
			o.c = s.Connect(o.Emited)

			b = append(b, o)
		}

		s.Emit(1)
		time.Sleep(2 * time.Second)
	}

	runtime.GC()

	time.Sleep(2 * time.Second)

	s.Emit(2)
	time.Sleep(2 * time.Second)
	runtime.GC()
	time.Sleep(2 * time.Second)
	s.Emit(3)
	fmt.Println("exit")
}
