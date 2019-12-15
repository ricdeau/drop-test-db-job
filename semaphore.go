package main

import "sync"

type Semaphore struct {
	counter int
	ch      chan bool
	once    sync.Once
}

func (s *Semaphore) Enter() {
	if s.counter == 0 {
		panic("Semaphore counter is 0")
	}
	if s.ch == nil {
		s.once.Do(func() {
			s.ch = make(chan bool, s.counter)
		})
	}
	s.ch <- true
}

func (s *Semaphore) Exit() {
	<-s.ch
}
