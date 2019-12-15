package main

import "sync"

type WaitingSemaphore struct {
	sync.WaitGroup
	counter int
	ch      chan bool
	once    sync.Once
}

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

func (s *WaitingSemaphore) Release() {
	<-s.ch
	s.Done()
}
