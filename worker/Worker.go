package worker

import (
	"sync"

	sync_mod "github.com/AnimusPEXUS/utils/sync"
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
type WaitExitResult chan EmptyStruct

type WorkerI interface {
	Start() WorkerControlChanResult
	Stop() WorkerControlChanResult
	Restart() WorkerControlChanResult
	Status() *workerstatus.WorkerStatusRO
	Wait() WaitExitResult
}

var _ WorkerI = &Worker{}

type Worker struct {
	status *workerstatus.WorkerStatus

	thread_func WorkerThreadFunction

	stop_flag        bool
	start_stop_mutex *sync.Mutex
	wait_lock        *sync_mod.MutexCheckable
	wait_cond        *sync.Cond

	//	signal_working *gosignal.Signal
	//	signal_stopped *gosignal.Signal
}

func New(f WorkerThreadFunction) *Worker {

	self := new(Worker)

	self.status = workerstatus.NewWorkerStatus(workerstatus.Stopped)

	self.thread_func = f
	self.start_stop_mutex = &sync.Mutex{}
	self.wait_lock = sync_mod.NewMutexCheckable(false)
	self.wait_cond = sync.NewCond(self.wait_lock)

	return self
}

func (self *Worker) Start() WorkerControlChanResult {
	ret := make(WorkerControlChanResult, 1)
	go func() {
		self.start_stop_mutex.Lock()
		defer self.start_stop_mutex.Unlock()

		if self.status.Stopped() {
			self.status.Set(workerstatus.Starting)
			self.stop_flag = false
			go func() {

				defer func() {
					self.stop_flag = true
					self.wait_cond.Broadcast()
					self.status.Reset()
				}()

				self.thread_func(
					func() {
						self.status.Set(workerstatus.Starting)
					},
					func() {
						self.status.Set(workerstatus.Working)
					},
					func() {
						self.status.Set(workerstatus.Stopping)
					},
					func() {
						self.status.Set(workerstatus.Stopped)
					},
					func() bool {
						return self.stop_flag
					},
				)
			}()
		} else {
			// TODO: probably, some error code should be reported in this case
			//       but this worker is intended to be working under some
			//       watchdog, which supposed to do it either way.
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

func (self *Worker) Status() *workerstatus.WorkerStatusRO {
	return &self.status.WorkerStatusRO
}

func (self *Worker) Wait() WaitExitResult {

	c := make(WaitExitResult, 1)
	if self.status.IsStopped() {
		c <- EmptyStruct{}
	} else {
		go func() {
			defer func() { c <- EmptyStruct{} }()
			self.wait_cond.Wait()
		}()
	}

	return c
}
