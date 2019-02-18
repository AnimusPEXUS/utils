package workerstatus

import "strings"

type WorkerStatus uint

const (
	Stopped WorkerStatus = iota
	Starting
	Working
	Stopping
)

func (self WorkerStatus) Stopped() bool {
	return self.IsStopped()
}

func (self WorkerStatus) IsStopped() bool {
	return self == Stopped
}

func (self *WorkerStatus) Reset() {
	*self = Stopped
	return
}

func (self *WorkerStatus) UpdateSelf(other *WorkerStatus) {
	*self = *other
}

func (self *WorkerStatus) UpdateOther(other *WorkerStatus) {
	*other = *self
}

func (self WorkerStatus) String() string {

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

func (self WorkerStatus) StringTitle() string {
	return strings.Title(self.String())
}

func (self WorkerStatus) StringT() string {
	return self.StringTitle()
}

func (self *WorkerStatus) Sum(in []WorkerStatus) {

	for _, i := range in {
		if i == Starting {
			*self = Starting
			return
		}
	}

	for _, i := range in {
		if i == Stopping {
			*self = Stopping
			return
		}
	}

	for _, i := range in {
		if i == Working {
			*self = Working
			return
		}
	}

	return

}
