package utils

import (
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func() error) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
