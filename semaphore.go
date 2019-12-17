package main

import "sync"

// WaitingSemaphore is a counting semaphore that can wait until all acquired lock will be released
type WaitingSemaphore struct {
	sync.WaitGroup
	counter int
	ch      chan bool
	once    sync.Once
}

// Acquire semaphore lock
func (s *WaitingSemaphore) Acquire() {
	if s.counter == 0 {
		panic("WaitingSemaphore counter is 0")
	}
	if s.ch == nil {
		s.once.Do(func() {
			s.ch = make(chan bool, s.counter)
		})
	}
	s.Add(1)
	s.ch <- true
}

// Release semaphore lock
func (s *WaitingSemaphore) Release() {
	<-s.ch
	s.Done()
}
