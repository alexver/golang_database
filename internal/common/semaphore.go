package common

import "sync"

type Semaphore struct {
	counter   int
	maxLimit  int
	condition *sync.Cond
}

func NewSemaphore(limit int) *Semaphore {
	mutex := &sync.Mutex{}

	return &Semaphore{
		maxLimit:  limit,
		condition: sync.NewCond(mutex),
	}
}

func (s *Semaphore) Acquire() {
	s.condition.L.Lock()
	defer s.condition.L.Unlock()

	for s.counter >= s.maxLimit {
		s.condition.Wait()
	}

	s.counter++
}

func (s *Semaphore) Release() {
	s.condition.L.Lock()
	defer s.condition.L.Unlock()

	s.counter--
	s.condition.Signal()
}
