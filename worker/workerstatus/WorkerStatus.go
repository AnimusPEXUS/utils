package workerstatus

import (
	"strings"
	"sync"
)

type WorkerStatusValue uint

const (
	Stopped WorkerStatusValue = iota
	Starting
	Working
	Stopping
)

func (self WorkerStatusValue) String() string {

	switch self {

	case Stopped:
		return "stopped"

	case Starting:
		return "starting"

	case Working:
		return "working"

	case Stopping:
		return "stopping"

	default:
		return "unknown"
	}

	return "error"
}

type WorkerStatus struct {
	WorkerStatusRO
}

type WorkerStatusRO struct {
	value WorkerStatusValue
	lock  *sync.RWMutex
}

func NewWorkerStatus(initial WorkerStatusValue) *WorkerStatus {
	self := &WorkerStatus{}
	self.value = initial
	self.lock = &sync.RWMutex{}
	return self
}

func (self *WorkerStatus) Set(value WorkerStatusValue) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.value = value
}

func (self *WorkerStatusRO) Get() WorkerStatusValue {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.value
}

func (self *WorkerStatusRO) Stopped() bool {
	return self.IsStopped()
}

func (self *WorkerStatusRO) IsStopped() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.Get() == Stopped
}

func (self *WorkerStatus) Reset() {
	self.Set(Stopped)
	return
}

func (self *WorkerStatus) UpdateSelf(other *WorkerStatusRO) {
	self.Set(other.Get())
}

func (self *WorkerStatusRO) UpdateOther(other *WorkerStatus) {
	other.Set(self.Get())
}

func (self *WorkerStatusRO) StringTitle() string {
	return strings.Title(self.Get().String())
}

func (self *WorkerStatusRO) StringT() string {
	return self.StringTitle()
}

func (self *WorkerStatus) Sum(in []*WorkerStatusRO) {

	for _, i := range in {
		if i.Get() == Starting {
			self.Set(Starting)
			return
		}
	}

	for _, i := range in {
		if i.Get() == Stopping {
			self.Set(Stopping)
			return
		}
	}

	for _, i := range in {
		if i.Get() == Working {
			self.Set(Working)
			return
		}
	}

	return

}
