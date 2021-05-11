package service

import (
	"fmt"
	"sync"
)

const bufferPoolSize = 32 * 1024

func newBufferPoolOld() *bufferPool {
	return &bufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				fmt.Println("ALLOC")
				return make([]byte, bufferPoolSize)
			},
		},
	}
}

func newBufferPool() *myAlloc {
	return &myAlloc{}
}

type myAlloc struct {
	pool [][]byte
	lock sync.Mutex
}

func (a *myAlloc) Get() []byte {
	var x []byte
	a.lock.Lock()
	if i := len(a.pool) - 1; i >= 0 {
		x = a.pool[i]
		a.pool = a.pool[:i]
	} else {
		x = make([]byte, bufferPoolSize)
	}
	a.lock.Unlock()
	return x
}
func (a *myAlloc) Put(x []byte) {
	a.lock.Lock()
	a.pool = append(a.pool, x)
	a.lock.Unlock()
}

type bufferPool struct {
	pool sync.Pool
}

func (b *bufferPool) Get() []byte {
	return b.pool.Get().([]byte)
}

func (b *bufferPool) Put(bytes []byte) {
	b.pool.Put(bytes)
}
