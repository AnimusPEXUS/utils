package worker

import (
	"sync"

	"github.com/AnimusPEXUS/utils/worker/workerstatus"
)

type WorkerThreadFunction func(

	set_starting func(),
	set_working func(),
	set_stopping func(),
	set_stopped func(),

	is_stop_flag func() bool,

)

type WorkerInterface interface {
	Start()
	Stop()
	Status() *workerstatus.WorkerStatus
}

type Worker struct {
	status *workerstatus.WorkerStatus

	thread_func WorkerThreadFunction

	stop_flag        bool
	start_stop_mutex *sync.Mutex

	//	signal_working *gosignal.Signal
	//	signal_stopped *gosignal.Signal
}

func New(f WorkerThreadFunction) *Worker {
	ret := new(Worker)

	ret.status = workerstatus.New()
	ret.thread_func = f
	ret.start_stop_mutex = &sync.Mutex{}

	return ret
}

func (self *Worker) Start() chan bool {
	ret := make(chan bool, 1)
	go func() {
		self.start_stop_mutex.Lock()
		defer self.start_stop_mutex.Unlock()

		if self.status.Stopped() {
			self.status.Starting = true
			self.stop_flag = false
			go func() {
				defer func() {
					self.stop_flag = true
					self.status.Reset()
				}()
				self.thread_func(
					func() {
						self.status.Starting = true
						self.status.Stopping = false
						self.status.Working = false
					},
					func() {
						self.status.Working = true
						self.status.Starting = false
						self.status.Stopping = false
					},
					func() {
						self.status.Stopping = true
						self.status.Starting = false
						self.status.Working = false
					},
					func() {
						self.status.Stopping = false
						self.status.Starting = false
						self.status.Working = false
					},
					func() bool {
						return self.stop_flag
					},
				)
			}()
		}
		ret <- true
	}()
	return ret
}

func (self *Worker) Stop() chan bool {
	ret := make(chan bool, 1)
	go func() {
		self.start_stop_mutex.Lock()
		defer self.start_stop_mutex.Unlock()

		self.stop_flag = true
		ret <- true
	}()
	return ret
}

func (self *Worker) Restart() chan bool {
	ret := make(chan bool, 1)
	go func() {
		<-self.Stop()
		<-self.Start()
		ret <- true
	}()
	return ret
}

func (self *Worker) Status() *workerstatus.WorkerStatus {
	return self.status
}
