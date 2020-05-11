package sync

import (
	sync_o "sync"
)

var _ LockerCheckable = &MutexCheckable{}

// difference to sync.Mutex, is what Unlock can be called if already unlocked
// this can be used in sync.Cond without afraiding to get error when Unlock() called
type MutexCheckable struct {
	mutex_o   *sync_o.Mutex
	is_locked bool
	s         *sync_o.Mutex
}

func NewMutexCheckable(locked bool) *MutexCheckable {
	self := &MutexCheckable{
		is_locked: locked,
		mutex_o:   &sync_o.Mutex{},
		s:         &sync_o.Mutex{},
	}
	if locked {
		self.mutex_o.Lock()
	}
	return self
}

func (self *MutexCheckable) Lock() {

	// to call or not to call self.s.Unlock() in defer
	var already_unlocked bool = false

	self.s.Lock()
	defer func() {
		if !already_unlocked {
			self.s.Unlock()
		}
	}()

	if self.is_locked {
		// to make object functions available
		already_unlocked = true
		go self.s.Unlock()
	}

	self.mutex_o.Lock()
	self.is_locked = true
}

func (self *MutexCheckable) Unlock() {

	self.s.Lock()
	defer self.s.Unlock()

	if !self.is_locked {
		return
	}

	self.mutex_o.Unlock()

	self.is_locked = false
}

func (self *MutexCheckable) IsLocked() bool {
	self.s.Lock()
	defer self.s.Unlock()
	return self.is_locked
}
