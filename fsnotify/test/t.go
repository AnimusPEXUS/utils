package main

import (
	"log"
	"time"

	"github.com/AnimusPEXUS/utils/fsnotify"
)

func main() {

	// log.Println(C.IN_CLOEXEC, fsnotify.IN_CLOEXEC, C.ccccc)

	// log.Println(C.IN_CLOSE, fsnotify.IN_CLOSE, C.ccccc2)
	// return

	w := fsnotify.NewWatcher(
		// fsnotify.IN_CLOEXEC|fsnotify.IN_NONBLOCK,
		fsnotify.IN_CLOEXEC,
		// 0,
		func(e *fsnotify.Event, err error) {
			if err != nil {
				log.Fatalln("event error:", err)
			}
			log.Println("event:", e)
		},
	)

	w2 := w.GetWorker()
	<-w2.Start()
	for {
		if w2.Status().StringT() == "Working" {
			break
		}
	}

	//fsnotify.IN_CLOSE_NOWRITE|fsnotify.IN_CLOSE_WRITE
	_, err := w.AddWatch("/home/animuspexus/tmp", fsnotify.IN_CLOSE)
	if err != nil {
		log.Fatalln("err add watch:", err)
	}

	for {
		time.Sleep(time.Second)
	}
}
