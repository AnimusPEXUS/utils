package workerstatus

import "strings"

type WorkerStatus struct {
	Starting bool
	Stopping bool
	Working  bool
}

func New() *WorkerStatus {
	ret := new(WorkerStatus)
	ret.Reset()
	return ret
}

func (self *WorkerStatus) Stopped() bool {
	return self.IsStopped()
}

func (self *WorkerStatus) IsStopped() bool {
	return !self.Starting && !self.Stopping && !self.Working
}

func (self *WorkerStatus) Reset() {
	self.Starting = false
	self.Stopping = false
	self.Working = false
	return
}

func (self *WorkerStatus) UpdateSelf(other *WorkerStatus) {
	self.Starting = other.Starting
	self.Stopping = other.Stopping
	self.Working = other.Working
}

func (self *WorkerStatus) UpdateOther(other *WorkerStatus) {
	other.Starting = self.Starting
	other.Stopping = self.Stopping
	other.Working = self.Working
}

func (self *WorkerStatus) String() string {

	if self.Starting && self.Stopping {
		return "invalid: starting and stopping"

	} else if self.Starting {
		return "starting"

	} else if self.Stopping {
		return "stopping"

	} else if self.Working {
		return "working"

	} else if self.IsStopped() {
		return "stopped"

	} else {
		return "unknown"
	}

	return "error"
}

func (self *WorkerStatus) StringTitle() string {
	return strings.Title(self.String())
}

func (self *WorkerStatus) StringT() string {
	return self.StringTitle()
}

func (self *WorkerStatus) Sum(in []*WorkerStatus) {

	res := New()

	for _, i := range in {
		if i.Starting {
			res.Starting = true
			goto exit
		}
	}

	for _, i := range in {
		if i.Stopping {
			res.Stopping = true
			goto exit
		}
	}

	for _, i := range in {
		if i.Working {
			res.Working = true
			goto exit
		}
	}

exit:
	self.UpdateSelf(res)

	return

}
