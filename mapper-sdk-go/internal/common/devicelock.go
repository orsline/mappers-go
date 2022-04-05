package common

import (
	"sync"
)

type Lock struct {
	DeviceLock *sync.Mutex
}

func (dl *Lock) Lock() {
	dl.DeviceLock.Lock()
}

func (dl *Lock) Unlock() {
	dl.DeviceLock.Unlock()
}
