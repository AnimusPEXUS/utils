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

type EmptyStruct struct{}

type WorkerControlChanResult chan EmptyStruct

type WorkerI interface {
	Start() WorkerControlChanResult
	Stop() WorkerControlChanResult
	Restart() WorkerControlChanResult
	Status() workerstatus.WorkerStatus
}

var _ WorkerI = &Worker{}

type Worker struct {
	status workerstatus.WorkerStatus

	thread_func WorkerThreadFunction

	stop_flag        bool
	start_stop_mutex *sync.Mutex

	//	signal_working *gosignal.Signal
	//	signal_stopped *gosignal.Signal
}

func New(f WorkerThreadFunction) *Worker {
	ret := new(Worker)

	ret.status = workerstatus.Stopped
	ret.thread_func = f
	ret.start_stop_mutex = &sync.Mutex{}

	return ret
}

func (self *Worker) Start() WorkerControlChanResult {
	ret := make(WorkerControlChanResult, 1)
	go func() {
		self.start_stop_mutex.Lock()
		defer self.start_stop_mutex.Unlock()

		if self.status.Stopped() {
			self.status = workerstatus.Starting
			self.stop_flag = false
			go func() {
				defer func() {
					self.stop_flag = true
					self.status.Reset()
				}()
				self.thread_func(
					func() {
						self.status = workerstatus.Starting
					},
					func() {
						self.status = workerstatus.Working
					},
					func() {
						self.status = workerstatus.Stopping
					},
					func() {
						self.status = workerstatus.Stopped
					},
					func() bool {
						return self.stop_flag
					},
				)
			}()
		}
		ret <- EmptyStruct{}
	}()
	return ret
}

func (self *Worker) Stop() WorkerControlChanResult {
	ret := make(WorkerControlChanResult, 1)
	go func() {
		self.start_stop_mutex.Lock()
		defer self.start_stop_mutex.Unlock()

		self.stop_flag = true
		ret <- EmptyStruct{}
	}()
	return ret
}

func (self *Worker) Restart() WorkerControlChanResult {
	ret := make(WorkerControlChanResult, 1)
	go func() {
		<-self.Stop()
		<-self.Start()
		ret <- EmptyStruct{}
	}()
	return ret
}

func (self *Worker) Status() workerstatus.WorkerStatus {
	return self.status
}
