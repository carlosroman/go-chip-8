package state

import (
	"io"
	"io/ioutil"
	"sync"
)

type Memory []uint8

func (m Memory) LoadMemory(r io.Reader) (err error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	for i := range b {
		m[i+512] = b[i]
	}
	return err
}

func InitMemory() Memory {
	return make(Memory, 4096)
}

type Stack struct {
	l sync.Mutex
	s []int16
	i int8
}

func (s *Stack) Pop() (val int16) {
	s.l.Lock()
	defer s.l.Unlock()
	val = s.s[s.i]
	s.i -= 1
	return val
}

func (s *Stack) Push(val int16) {
	s.l.Lock()
	defer s.l.Unlock()
	s.i += 1
	s.s[s.i] = val
}

func (s *Stack) Len() (length int8) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.i + 1
}

func InitStack() *Stack {
	return &Stack{
		s: make([]int16, 16),
		i: -1,
	}
}
