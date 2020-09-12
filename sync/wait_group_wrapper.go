package sync

import (
	gosync "sync"
)

type WaitGroupWrapper struct {
	gosync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
